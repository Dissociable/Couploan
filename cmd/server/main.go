package main

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/ent/proxy"
	"github.com/Dissociable/Couploan/ent/proxyprovider"
	"github.com/Dissociable/Couploan/ent/user"
	"github.com/Dissociable/Couploan/pkg/routes"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/Dissociable/Couploan/pkg/tasks"
	"github.com/Dissociable/Couploan/proxstore"
	"github.com/Dissociable/Couploan/ve"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		Shutdown()
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()
	if err := run(ctx, os.Args); err != nil {
		c.Logger.Error("run returned with error", zap.Error(err))
		Shutdown()
		os.Exit(exitCodeErr)
	}
}

func run(ctx context.Context, _ []string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := runMain(ctx)
			return err
		}
	}
}

// runMain starts the bot
func runMain(ctx context.Context) (err error) {
	// Start a new container
	c = services.NewContainer()
	if c.Config.App.Environment == config.EnvLocal || c.Config.App.Environment == config.EnvDevelop {
		err = prepareForDevRun(ctx)
		if err != nil {
			return err
		}
	}

	err = loadProxies(ctx, c)
	if err != nil {
		c.Logger.Error("failed to load proxies", zap.Error(err))
		return
	}

	c.Logger.Info("Loaded proxies", zap.Int("count", c.ProxyStore.Count()))

	if c.Config.App.Environment == config.EnvLocal || c.Config.App.Environment == config.EnvDevelop {
		p := c.ProxyStore.Next()
		v := ve.New(c.ProxyStore, p)
		// for i := 0; i < 6; i++ {
		// 	ip, err := v.IP(ctx)
		// 	if err != nil {
		// 		c.Logger.Error("failed to get IP", zap.Error(err))
		// 		return err
		// 	}
		// 	c.Logger.Info("IP", zap.String("ip", ip))
		// 	released, err := c.ProxyStore.ReleaseProxy(p)
		// 	if err != nil {
		// 		c.Logger.Error("failed to release proxy", zap.Error(err))
		// 		return err
		// 	}
		// 	if !released {
		// 		c.Logger.Fatal("failed to released proxy", zap.String("proxy", p.String()))
		// 	}
		// }
		ip, err := v.IP(ctx)
		if err != nil {
			c.Logger.Error("failed to get IP", zap.Error(err))
			return err
		}
		c.Logger.Info("IP", zap.String("ip", ip))
		_, err = v.Index(ctx)
		if err != nil {
			c.Logger.Error("failed to get index", zap.Error(err))
			return err
		} else {
			c.Logger.Debug("got index successfully")
		}
	}

	// Start the scheduler service to queue periodic tasks
	tasks.StartTasksRunner(c, true)

	// Start the bot
	routes.BuildRouter(c)

	return c.StartServer(ctx)
}

func loadProxies(ctx context.Context, container *services.Container) (err error) {
	providers := make(map[int]*proxstore.Provider)
	proxies, err := c.ORM.Proxy.Query().WithProxyProvider().All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = nil
			return
		}
		err = errors.Wrap(err, "failed to load proxies")
		return
	}
	for _, p := range proxies {
		var prox *proxstore.Proxy[tls_client.HttpClient]
		if p.Username != nil && p.Password != nil {
			prox = proxstore.NewProxyWithCredential[tls_client.HttpClient](
				p.IP,
				p.Port,
				proxstore.Protocol(strings.ToLower(string(p.Type))),
				*p.Username,
				*p.Password,
			)
		} else {
			prox = proxstore.NewProxy[tls_client.HttpClient](
				p.IP,
				p.Port,
				proxstore.Protocol(strings.ToLower(string(p.Type))),
			)
		}
		if p.Edges.ProxyProvider != nil {
			existingProvider, ok := providers[p.Edges.ProxyProvider.ID]
			if !ok {
				existingProvider = &proxstore.Provider{
					Name:        proxstore.ProviderName(p.Edges.ProxyProvider.Name),
					ServiceType: p.Edges.ProxyProvider.ServiceType,
					Username:    p.Edges.ProxyProvider.Username,
					Password:    p.Edges.ProxyProvider.Password,
				}
				providers[p.Edges.ProxyProvider.ID] = existingProvider
			}
			prox = prox.SetProvider(existingProvider)
		}
		if p.Rotating {
			prox.Rotating = true
		}

		err = container.ProxyStore.LoadProxy(prox)
		if err != nil {
			err = errors.Wrapf(err, "failed to load proxy: %s", prox.String())
			return
		}
	}

	return nil
}

// prepareForDevRun sets up dev environment
func prepareForDevRun(ctx context.Context) (err error) {
	err = c.ORM.User.Create().
		SetName(gofakeit.Name()).
		SetKey(strings.Repeat("test", 32/4)).
		SetBalance(100).
		SetContact("Telegram: "+gofakeit.Username()).
		SetRole(user.RoleAdmin).
		OnConflict(
			sql.ConflictColumns(user.FieldKey),
			sql.ResolveWithNewValues(),
		).
		Exec(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to create user for dev environment")
		return
	}

	_, err = c.ORM.Proxy.Delete().Exec(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to delete proxies for dev environment")
		return
	}
	if len(c.Config.Tests.Proxy.Lines) > 0 {
		provider := proxstore.Provider{
			Name:        proxstore.ProviderName(c.Config.Tests.Proxy.Provider.Name),
			ServiceType: c.Config.Tests.Proxy.Provider.Service,
			Username:    c.Config.Tests.Proxy.Provider.Username,
			Password:    c.Config.Tests.Proxy.Provider.Password,
		}
		providerId := 0
		if provider.Name != proxstore.ProviderNameNone {
			providerId, err = c.ORM.ProxyProvider.Create().
				SetName(string(provider.Name)).
				SetServiceType(provider.ServiceType).
				SetUsername(provider.Username).
				SetPassword(provider.Password).
				OnConflict(
					sql.ConflictColumns(
						proxyprovider.FieldName,
						proxyprovider.FieldServiceType,
						proxyprovider.FieldUsername,
						proxyprovider.FieldPassword,
					),
					sql.ResolveWithNewValues(),
				).ID(ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to create proxy provider for dev environment")
				return err
			}
		}
		for _, line := range c.Config.Tests.Proxy.Lines {
			p, err := proxstore.ParseLineWithoutProtocol[tls_client.HttpClient](
				line,
				strings.Split(line, ":"),
				proxstore.ProtocolSocks5,
			)
			if err != nil {
				err = errors.Wrap(err, "failed to parse proxy for dev environment")
				return err
			}
			if p == nil {
				err = errors.New("failed to parse proxy for dev environment")
				return err
			}
			q := c.ORM.Proxy.Create().
				SetType(proxy.Type(strings.ToUpper(string(p.Protocol)))).
				SetIP(p.Host).
				SetPort(p.Port).
				SetRotating(true).
				SetUsername(p.Username).
				SetPassword(p.Password)
			if providerId > 0 {
				q = q.SetProxyProviderID(providerId)
			}
			err = q.OnConflict(
				sql.ConflictColumns(proxy.FieldIP, proxy.FieldPort, proxy.FieldUsername, proxy.FieldPassword),
				sql.ResolveWithNewValues(),
			).
				Exec(ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to create user for dev environment")
				return err
			}
		}
	}

	return nil
}

func Shutdown() {
	if c != nil && c.Logger != nil {
		c.Logger.Info("Exiting...")
	} else {
		fmt.Println("Exiting...")
	}
	if c.Web != nil {
		if err := c.Web.Shutdown(); err != nil {
			c.Logger.Error(err.Error())
		}
	}
	if c != nil {
		if err := c.Shutdown(); err != nil {
			fmt.Println(err.Error())
		}
	}
}

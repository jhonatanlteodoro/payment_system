package shared_deps

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var deps *SharedDeps
var once sync.Once

type envs struct {
	DbUser string `env:"DB_USER" envDefault:"secret_user"`
	DbPass string `env:"DB_PASS" envDefault:"secret_password"`
	DbHost string `env:"DB_HOST" envDefault:"localhost"`
	DbPort string `env:"DB_PORT" envDefault:"5432"`
	DbName string `env:"DB_NAME" envDefault:"payment"`

	RabbitMqUser string `env:"RABBITMQ_USER" envDefault:"secret_user"`
	RabbitMqPass string `env:"RABBITMQ_PASS" envDefault:"secret_password"`
	RabbitMqHost string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	RabbitMqPort string `env:"RABBITMQ_PORT" envDefault:"5672"`

	StartPaymentQueueName   string `env:"START_PAYMENT_QUEUE_NAME" envDefault:"start-payment"`
	ProcessPaymentQueueName string `env:"PROCESS_PAYMENT_QUEUE_NAME" envDefault:"process-payment"`
	NotifyUserQueueName     string `env:"NOTIFY_USER_QUEUE_NAME" envDefault:"notify-user"`
}

// SharedDeps are intended to hold env vars and shared items [db connections, clients, and etc...]
type SharedDeps struct {
	Envs *envs

	PaymentsDBConn *pgx.Conn
	RabbitMqConn   *amqp.Connection

	//Queues
	StartPaymentQueue   ports.Queue
	ProcessPaymentQueue ports.Queue
	NotifyUserQueue     ports.Queue
}

// NewSharedDependencies will return a new instance of dependencies and register background watcher so we can close any required object
// This must be called once in the application start. For random settings access use GetSharedDependencies function
func NewSharedDependencies(shutdown chan os.Signal) *SharedDeps {
	once.Do(func() {
		e := &envs{}
		if err := env.Parse(e); err != nil {
			log.Fatalf("Failed parsing envs %v", err)
		}

		deps = &SharedDeps{
			Envs:           e,
			PaymentsDBConn: getPaymentsDBConn(e),
			RabbitMqConn:   getRabbitMqConn(e),
		}
		deps.StartPaymentQueue = NewQueue(deps.RabbitMqConn, e.StartPaymentQueueName)
		deps.ProcessPaymentQueue = NewQueue(deps.RabbitMqConn, e.ProcessPaymentQueueName)
		deps.NotifyUserQueue = NewQueue(deps.RabbitMqConn, e.NotifyUserQueueName)

		go func() {
			<-shutdown
			// register here any item that must be properly closed when shutting down
			log.Println("Shared Dependencies Shutting down...")
			ctx := context.Background()
			deps.RabbitMqConn.Close()
			log.Println("Rabbit MQ connection closed...")
			deps.PaymentsDBConn.Close(ctx)
			log.Println("Payments DB connection closed...")
			shutdown <- syscall.SIGINT // propagate
		}()
	})
	return deps
}

// GetSharedDependencies returns the current settings
// this method will not initialize the settings, if you need it call NewSharedDependencies in the application start instead
func GetSharedDependencies() *SharedDeps {
	return deps
}

func getPaymentsDBConn(env *envs) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", env.DbUser, env.DbPass, env.DbHost, env.DbPort, env.DbName)
	dbConn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if errPing := dbConn.Ping(ctx); errPing != nil {
		log.Fatalf("DB Ping failed: %v\n", errPing)
	}

	log.Println("Payment DB connection established")
	return dbConn
}

func getRabbitMqConn(env *envs) *amqp.Connection {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", env.RabbitMqUser, env.RabbitMqPass, env.RabbitMqHost, env.RabbitMqPort))
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v\n", err)
	}

	log.Println("RabbitMQ connection established")
	return conn
}

package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log/slog"
	"os"
	"product-tracking/config"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaSvc struct {
	Settings *config.Settings
}

func (k KafkaSvc) SendMessage(msg string, key string, topic string) {
	// Create a cancellable context for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture termination signals
	/*sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Info("Shutting down application...")
		cancel()
	}()*/

	// Configure SASL Mechanism (SCRAM-SHA-512 is common)
	mechanism, err := scram.Mechanism(scram.SHA512, k.Settings.K_USER, k.Settings.K_PASS)
	if err != nil {
		slog.Error("failed to create SASL mechanism", "description", err.Error())
	}

	// Configure TLS / SSL
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Using self-signed certificates, load CA cert:
	caCert, err := os.ReadFile("ca.pem")
	if err == nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	// Create a Shared Transport with SASL + SSL
	sharedTransport := &kafka.Transport{
		SASL: mechanism,
		TLS:  tlsConfig,
	}

	// Start Producer and Consumer concurrently
	k.sendMsgToKafka(ctx, sharedTransport, msg, key, topic)
}

func (k KafkaSvc) sendMsgToKafka(ctx context.Context, transport *kafka.Transport, msg string, key string, topic string) {
	slog.Info("Sending message, connecting to kafka..")
	writer := &kafka.Writer{
		Addr:  kafka.TCP(k.Settings.K_BROKER_ADDRESS),
		Topic: topic,
		//Balancer:     &kafka.LeastBytes{}, // Distributes messages across partitions
		Async: false, // Synchronous writes for strict error checking
		//RequiredAcks: kafka.RequireAll,    // Ensures acknowledgment from all synced replicas
		Transport: transport,
	}
	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(msg),
	})
	if err != nil {
		slog.Error("Producer failed to write message", "description", err.Error())
	} else {
		slog.Info("Kafka message sended.")
	}

}

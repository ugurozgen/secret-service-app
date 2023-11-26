package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

// export VAULT_ADDR='http://127.0.0.1:8200'
// export VAULT_TOKEN=<VAULT_TOKEN>
func vaultClient() *vault.Client {
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(os.Getenv("VAULT_ADDR")),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(os.Getenv("VAULT_TOKEN")); err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	client := vaultClient()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/secret", func(c *gin.Context) {
		var newSecret map[string]any
		
		err := c.Bind(&newSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err = client.Secrets.KvV2Write(context.Background(), "app", schema.KvV2WriteRequest{Data: newSecret}, vault.WithMountPath("kv-v2"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "secret created",
		})
	})

	r.GET("/secret/:secret", func(c *gin.Context) {
		secret := c.Param("secret")

		s, err := client.Secrets.KvV2Read(context.Background(), secret, vault.WithMountPath("kv"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, s.Data.Data)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

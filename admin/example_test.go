package admin_test

import (
	"context"
	"fmt"
	"log"

	"github.com/zero-color/ecforce-go"
	"github.com/zero-color/ecforce-go/admin"
)

func ExampleNewClient() {
	ctx := context.Background()

	client, err := admin.NewClient("https://example.ec-force.com",
		ecforce.WithToken("LcPTYtBqpUFQag3GwqAVFyv_EEcsFvEC"))
	if err != nil {
		log.Fatal(err)
	}

	opts := &ecforce.ListOptions{
		Q:    ecforce.Query{"email_cont": "@example.com"},
		Sort: []string{"-created_at"},
		Per:  100,
	}
	for {
		customers, resp, err := client.Customers.List(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range customers {
			fmt.Println(c.Attributes.Email)
		}
		if !resp.HasNextPage() {
			break
		}
		opts.Page = resp.NextPage()
	}
}

func ExampleCustomersService_Create() {
	ctx := context.Background()

	client, err := admin.NewClient("https://example.ec-force.com",
		ecforce.WithToken("LcPTYtBqpUFQag3GwqAVFyv_EEcsFvEC"))
	if err != nil {
		log.Fatal(err)
	}

	created, _, err := client.Customers.Create(ctx, &admin.CustomerCreateRequest{
		Customer: &admin.CustomerParams{
			Email:    ecforce.String("test@example.com"),
			Password: ecforce.String("secret-password"),
			State:    ecforce.String("member"),
			Optin:    ecforce.Bool01(true),
			BillingAddressAttributes: &admin.AddressParams{
				Name01:       ecforce.String("テスト"),
				Name02:       ecforce.String("太郎"),
				Kana01:       ecforce.String("テスト"),
				Kana02:       ecforce.String("タロウ"),
				Zip01:        ecforce.String("111"),
				Zip02:        ecforce.String("1111"),
				PrefectureID: ecforce.Int64(13),
				Addr01:       ecforce.String("目黒区下目黒"),
				Addr02:       ecforce.String("2-23-18"),
				Tel01:        ecforce.String("03"),
				Tel02:        ecforce.String("5759"),
				Tel03:        ecforce.String("6380"),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(created.Attributes.Number)
}

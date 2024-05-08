package prisma

import "github.com/sriramk1n1/C_Platform_Backend/prisma/db"

func HandleDBOperation(operation func(client *db.PrismaClient) error) error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	return operation(client)
}

package prisma

import "first.com/prisma/db"

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

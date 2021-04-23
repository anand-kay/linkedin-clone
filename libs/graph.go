package libs

import (
	"context"
	"database/sql"
	"log"

	"github.com/anand-kay/linkedin-clone/models"
	"github.com/go-redis/redis/v8"
)

// PopulateGraph - Populates connections graph at startup
func PopulateGraph(ctx context.Context, db *sql.DB, rdb *redis.Client) error {
	connections, err := models.GetConnections(db)
	if err != nil {
		return err
	}

	for _, connection := range connections {
		err := connection.AddToGraph(ctx, rdb)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckLevel - Checks connection level between two IDs
func CheckLevel(ctx context.Context, rdb *redis.Client, rootID string, searchID string) uint8 {
	if rootID == searchID {
		return 0
	}

	var idQueue []string
	touchedMap := make(map[string]bool)

	if isLevelOne(ctx, rdb, rootID, searchID, &idQueue, &touchedMap) {
		return 1
	}

	if isLevelTwo(ctx, rdb, rootID, searchID, &idQueue, &touchedMap) {
		return 2
	}

	if isLevelThree(ctx, rdb, rootID, searchID, &idQueue, &touchedMap) {
		return 3
	}

	return 99
}

func isLevelOne(ctx context.Context, rdb *redis.Client, rootID string, searchID string, idQueue *[]string, touchedMap *map[string]bool) bool {
	ok, err := rdb.SIsMember(ctx, "user:"+rootID, searchID).Result()
	if err != nil {
		log.Println(err)

		return false
	}

	if ok {
		return true
	}

	ids, err := rdb.SMembers(ctx, "user:"+rootID).Result()
	if err != nil {
		log.Println(err)

		return false
	}

	for _, id := range ids {
		*idQueue = append(*idQueue, id)
		(*touchedMap)[id] = true
	}

	return false
}

func isLevelTwo(ctx context.Context, rdb *redis.Client, rootID string, searchID string, idQueue *[]string, touchedMap *map[string]bool) bool {
	for _, id := range *idQueue {
		ok, err := rdb.SIsMember(ctx, "user:"+id, searchID).Result()
		if err != nil {
			log.Println(err)

			return false
		}

		if ok {
			return true
		}

		ids, err := rdb.SMembers(ctx, "user:"+id).Result()
		if err != nil {
			log.Println(err)

			return false
		}

		for _, id := range ids {
			if _, ok := (*touchedMap)[id]; !ok {
				*idQueue = append(*idQueue, id)
				(*touchedMap)[id] = true
			}
		}

		*idQueue = (*idQueue)[1:]
	}

	return false
}

func isLevelThree(ctx context.Context, rdb *redis.Client, rootID string, searchID string, idQueue *[]string, touchedMap *map[string]bool) bool {
	for _, id := range *idQueue {
		ok, err := rdb.SIsMember(ctx, "user:"+id, searchID).Result()
		if err != nil {
			log.Println(err)

			return false
		}

		if ok {
			return true
		}
	}

	return false
}

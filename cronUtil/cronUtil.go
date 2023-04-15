package cronUtil

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"github.com/windbnb/user-service/client"
	model "github.com/windbnb/user-service/model"
)

func ConfigureCronJobs(db *gorm.DB) {
	cronHandler := cron.New()
	cronHandler.AddFunc("@hourly", func() {
		var userDeletionRequests []model.UserDeletionEvent
		db.Find(&userDeletionRequests)

		for _, userDeletionRequest := range userDeletionRequests {
			userId := userDeletionRequest.UserId
			fmt.Println(userId)
			err := client.DeleteAccomodationForHost(uint(userDeletionRequest.UserId))
			if err == nil {
				db.Delete(userDeletionRequest)
			}
		}
	})
}

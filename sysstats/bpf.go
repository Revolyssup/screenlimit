package sysstats

import "gorm.io/gorm"

type StatCollector struct {
	programs []string //eg- brave, chrome, firefox or any other app we want to monitor.
	*gorm.DB
}

//TODO: It will run in the background extracting statistics and pushing them into the database
func (s *StatCollector) Run() {

}

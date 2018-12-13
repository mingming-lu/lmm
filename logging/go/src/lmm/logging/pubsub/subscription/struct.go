package subscription

type ZapBase struct {
	Level      string `datastore:"level"   json:"level"`
	TimeStamp  string `datastore:"time"    json:"ts"`
	LoggerName string `datastore:"-"       json:"logger"`
	Message    string `datastore:"message" json:"msg"`
}

package task

type Task struct {
	ID           string  `db:"id" json:"id"`
	TrackerID    string  `db:"tracker_id" json:"trackerId"`
	Interval     float64 `db:"interval" json:"interval"`
	IntervalUnit string  `db:"interval_unit" json:"intervalUnit"`
	Time         string  `db:"time" json:"time"`
	Remark       string  `db:"remark" json:"remark"`
	Created      string  `db:"created_at" json:"created_at"`
	Updated      string  `db:"updated_at" json:"updated_at"`
}

func All() {

}

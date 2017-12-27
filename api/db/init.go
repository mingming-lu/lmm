package db

import "log"

func Init() {
	d := New().Create("lmm").Use("lmm")
	defer d.Close()

	for _, sql := range CreateSQL {
		_, err := d.Exec(sql)
		if err != nil {
			log.Println(err)
		}
	}
}

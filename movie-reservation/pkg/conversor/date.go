package conversor

import "time"

var layoutMap = map[string]string{
    "dateLayout":     "02/01/2006",
    "dateTimeLayout": "2006-01-02 15:04:05",
}


func StrToDateTime(strDate, layout string) (time.Time, error) {
	date, err := time.Parse(layoutMap[layout], strDate)
	if err != nil{
		return time.Time{}, err
	}
	return date, nil
}

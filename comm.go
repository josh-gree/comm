package comm
// A comment
import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/labstack/echo"
)

type Message interface{
	Recieve()
	Send()
}

type JobMessage struct{
	Id int `json:"id"`
	Data []float64 `json:"data"`
	Service string `json:"service"`
}

var locs = map[string]string{"sum":"sum:8000/job","prod":"prod:8000/job"}

func (j *JobMessage) Recieve(public bool, service ...func(data []float64, id int)) func(c echo.Context) error{
	return func (c echo.Context) error{
		err := c.Bind(&j)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("%#v\n",j)
		if public {
			// Send Job to service
			j.Send(locs[j.Service])
		} else {
			// Do Job
			// Send result to public
			go service[0](j.Data,j.Id)
		}

		return nil
	}
}

func (j *JobMessage) Send(dest string) error{
	data, err := json.Marshal(j)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = http.Post(fmt.Sprintf("http://%s",dest),"application/json",bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type ResMessage struct{
	Id int `json:"id"`
	Result float64 `json:"result"`
}

func (r *ResMessage) Recieve(public bool) func(c echo.Context) error{
	return func (c echo.Context) error{
		err := c.Bind(&r)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if public {
			// Log result to stdout
			fmt.Printf("%#v\n",r)
		} else {
			// This will never happen (services do not recieve results)
		}

		return nil
	}
}

func (r *ResMessage) Send(dest string) error{
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = http.Post(fmt.Sprintf("http://%s",dest),"application/json",bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

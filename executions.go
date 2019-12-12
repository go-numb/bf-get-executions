package eoy // Executions of year

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/go-numb/go-bitflyer/auth"
	"github.com/go-numb/go-bitflyer/v1"

	"github.com/labstack/gommon/log"

	"github.com/go-numb/go-bitflyer/v1/private/executions"

	"github.com/go-numb/go-bitflyer/v1/types"
)

// New creates Exchange struct by name
func New(key, secret, productCode string) *Bitflyer {
	return &Bitflyer{
		C: v1.NewClient(&v1.ClientOpts{
			AuthConfig: &auth.AuthConfig{
				APIKey:    key,
				APISecret: secret,
			},
		}),
		Code: types.ProductCode(productCode),
	}
}

// Bitflyer includes Exchange interface
type Bitflyer struct {
	C    *v1.Client
	Code types.ProductCode
}

// Executions gets executions
func (p *Bitflyer) Executions() {
	var (
		lastID, lastExecID int
		lastYear           int
	)

	f, err := os.OpenFile(
		fmt.Sprintf("./%s_executions_bitflyer.csv", time.Now().Format("20060102150405")),
		os.O_CREATE|os.O_WRONLY,
		0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{"index", "id", "accept_id", "child_id", "side", "price", "size", "commission", "created_at"})

	for { // 1reqest/10sec = 30/5minutes
		time.Sleep(10 * time.Second)
		exec, _, err := p.C.ExecutionsMe(&executions.Request{
			ProductCode: p.Code,
			Pagination: types.Pagination{
				Count:  499,
				Before: lastExecID,
			},
		})
		if err != nil {
			log.Error(err)
			continue
		}

		for i, ex := range *exec {
			w.Write([]string{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%d", ex.ID),
				ex.ChildOrderAcceptanceID,
				ex.ChildOrderID,
				ex.Side,
				fmt.Sprintf("%.f", ex.Price),
				fmt.Sprintf("%.12f", ex.Size),
				fmt.Sprintf("%.f", ex.Commission),
				ex.ExecDate.Time.String(),
			})

			lastID++
		}
		w.Flush()

		var execs []executions.Execution
		execs = *exec
		l := len(execs) - 1
		if l < 1 {
			break
		}
		lastExecID = execs[l].ID

		progress := execs[l].ExecDate.Time
		fmt.Printf("progress... %d - %+v\n", lastID, progress.Format("2016/01/02 15:04"))
		lastYear = progress.Year()
		if lastYear != 2019 {
			break
		}
	}

	fmt.Printf("done...%+v\n", lastExecID)
}

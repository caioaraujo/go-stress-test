package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-stress-test",
	Short: "Teste de carga em servidor WEB",
	Long: `Aplicação de teste de carga em servidor WEB. Parâmetros:
	--url: exemplo, http://www.google.com ;
	--requests: total de requisições. Exemplo, 100;
	--concurrency: quantidade de concorrências. Exemplo: 5.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			cmd.PrintErrln("Erro na URL:", err)
			return
		}
		if url == "" {
			cmd.PrintErrln("URL é obrigatório! Favor fornecer um valor válido.")
			return
		}

		requests, err := cmd.Flags().GetUint("requests")
		if err != nil {
			cmd.PrintErrln("Erro no argumento requests:", err)
			return
		}
		if requests == 0 {
			cmd.PrintErrln("Requests deve ser maior que zero.")
			return
		}

		concurrency, err := cmd.Flags().GetUint("concurrency")
		if err != nil {
			cmd.PrintErrln("Erro no argumento concurrency:", err)
			return
		}
		if concurrency == 0 {
			cmd.PrintErrln("Concurrency deve ser maior que zero.")
			return
		}

		if concurrency > requests {
			cmd.PrintErrln("Quantidade de requests deve ser igual ou superior ao total de concorrências")
			return
		}

		reqTotalRealizadas := 0
		req200Total := 0
		reqOutrosStatusTotal := make(map[int]int)
		tempoInicio := time.Now()

		var wg sync.WaitGroup
		wg.Add(int(concurrency))
		reqChan := make(chan int, requests)

		for i := 0; i < int(concurrency); i++ {
			go func(i int, ci <-chan int) {
				defer wg.Done()
				j := 1
				for req := range ci {
					fmt.Println("Request nº", req)
					resp, err := http.Get(url)
					if err != nil {
						cmd.PrintErrln("Erro ao realizar requisição:", err)
						os.Exit(1)
					}
					reqTotalRealizadas++
					if resp.StatusCode == http.StatusOK {
						req200Total++
					} else {
						reqOutrosStatusTotal[resp.StatusCode]++
					}
					j++
				}
			}(i, reqChan)
		}

		for j := 0; j < int(requests); j++ {
			reqChan <- j
		}
		close(reqChan)
		wg.Wait()

		duracao := time.Since(tempoInicio)

		cmd.Println("Tempo total gasto na execução:", duracao)
		cmd.Println("Requests realizadas:", reqTotalRealizadas)
		cmd.Println("Requests status HTTP 200 realizadas:", req200Total)
		cmd.Println("Requests em outros status:", reqOutrosStatusTotal)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "URL para teste")
	rootCmd.Flags().UintP("requests", "r", 0, "Total de requisições")
	rootCmd.Flags().UintP("concurrency", "c", 0, "Total de concorrências")
}

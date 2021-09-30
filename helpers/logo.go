package helpers

import "fmt"

var colorReset = "\033[0m"

var (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var arrow string = "➤"

var logo string = "            .......           \n     .,lxOkc;,,,,;oK0xc.      \n   ;kkc'xX,.,:ldxdoc,',dO;    \n .Oo.  lX0NNOl,.       .ONk.  \n lK  c0l. oK.     .'cxkdkOo0  \n  xKXk.    KXo:lx0Kx;.  ol 0c \n   .cdolcldXNNNKo'      O; kl \n          .;OXxoc;.    .X'.K' \n             c0: .;lo:'xN;kc  \n              .ldlcclkXNXx.   \n                     ...   "

func PrintLogo(url string) {
	fmt.Println(logo)
	fmt.Printf(" %s http server started on %s\n", arrow, string(colorGreen)+url+string(colorReset))
}

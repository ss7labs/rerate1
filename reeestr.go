package main
import (
 "os"
 "fmt"
 "sync"
 "strconv"
)
type Total struct {
 totalInt float64
 totalSng float64
 totalTkm float64
 totalLoc float64
 totalPstn float64
 total09 float64
 totalTen float64
 totalMan float64
}
type ReestrGroup struct {
 file *os.File
 total Total
}

type Reestr struct {
 r709 ReestrGroup
}

func (r *Reestr) createReestrFiles(date string) {
//R_709_20.O07
 var err error
 runes := []rune(date)
 YY := string(runes[2:4])
 MM := string(runes[4:])
 fn := "R_709_"+YY+".O"+MM
 r.r709 = ReestrGroup{}
 r.r709.file,err = os.Create(fn)
 check(err)
 r.r709.file.WriteString("     Реестp по оpганизациям за август 2020 года. \n")
 r.r709.file.WriteString("     ( оплата в манатах)      \n")
 r.r709.file.WriteString("!----------------------------------------------------------------------------------------------------------------------------------------------------!\n")
 r.r709.file.WriteString("!  Л/счет  !    НАЗВАНИЕ ОРГАНИЗАЦИИ        !                   Пеpеговоpы ( сумма в манатах )                                                       !\n")
 r.r709.file.WriteString("!          !                                !Междунаpодн.!   По СНГ   !Туpкменистан! Городск.   !Услуги ATC  !Услуги 09   !Услуги(10%) !   Итого     !\n")
 r.r709.file.WriteString("!----------------------------------------------------------------------------------------------------------------------------------------------------!\n")

 fmt.Println("createReestrFiles",fn)
}
func (r *Reestr) closeFiles() {
 r.r709Total()
 r.r709.file.Close()
}
func (r *Reestr) r709Total() {
 r.r709.file.WriteString("!----------------------------------------------------------------------------------------------------------------------------------------------------!\n")
 //!     И Т О Г О                            !      60.12 !       0.00 !     574.56 !     2.41 !    51.60 !    30.00 !     20.58 !      739.27 !
 str := "!     И Т О Г О                             "
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalInt),12)
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalSng),12)
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalTkm),12)
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalLoc),12)
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalPstn),12)
 str += alignCenter("0.00",12)
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalTen),12)
 str += alignCenter(fmt.Sprintf("%.2f",r.r709.total.totalMan),13)
 str += "!\n"
 r.r709.file.WriteString(str)
}
/*
     Реестp по оpганизациям за июль 2020 года. 
     ( оплата в манатах)      
!--------------------------------------------------------------------------------------------------------------------------------------------!
!  Л/счет  !    НАЗВАНИЕ ОРГАНИЗАЦИИ       !                   Пеpеговоpы ( сумма в манатах )                                                !
!          !                               !Междунаpодн.!   По СНГ   !Туpкменистан! Городск. !Услуги ATC!Услуги 09 !Услуги(10%)!   Итого     !
!--------------------------------------------------------------------------------------------------------------------------------------------!
! 70984500 ! Довлетабадгазчыкарыш Дервезин !       0.00 !       0.00 !       0.00 !     0.29 !     0.00 !     0.00 !      0.00 !        0.29 !
!--------------------------------------------------------------------------------------------------------------------------------------------!
!     И Т О Г О                            !      60.12 !       0.00 !     574.56 !     2.41 !    51.60 !    30.00 !     20.58 !      739.27 !
*/
var fileMutex sync.Mutex
func (org *Org) addToReestr() {
 //str := "! 70984500 ! Довлетабадгазчыкарыш Дервезин !       0.00 !       0.00 !       0.00 !     0.29 !     0.00 !     0.00 !      0.00 !        0.29 !\n"
 intTotal:=fmt.Sprintf("%.2f",org.IntCallsTotal.Man)
 sngTotal:=fmt.Sprintf("%.2f", org.SngCallsTotal.Man)
 tkmTotal:=fmt.Sprintf("%.2f", org.TkmCallsTotal.Man)
 locTotal:=fmt.Sprintf("%.2f", org.LocalCallsTotal.Man)
 pstnTotal:=fmt.Sprintf("%.2f", org.PstnServiceTotal.Man)

 tenTotal:=fmt.Sprintf("%.2f", org.tenPercentMan)
 orgTotal:=fmt.Sprintf("%.2f",org.orgTotalMan)
 //Field 1
 str:= alignCenter(strconv.Itoa(org.Acct),10)
 str += alignLeft(org.Name,31)
 str += alignCenter(intTotal,12)
 str += alignCenter(sngTotal,12)
 str += alignCenter(tkmTotal,12)
 str += alignCenter(locTotal,12)
 str += alignCenter(pstnTotal,12)
 str += alignCenter("0.00",12)
 str += alignCenter(tenTotal,12)
 str += alignCenter(orgTotal,13)
 str += "!\n"
 //str := fmt.Sprintf("!%10s!%31s!%12s!%12s!%12s!%12s!%12s!%12s!%12s!%13s!",strconv.Itoa(org.Acct),org.Name,intTotal,sngTotal,tkmTotal,locTotal,pstnTotal,"0.00",tenTotal,orgTotal)
 org.reestr.r709.total.totalInt += org.IntCallsTotal.Man
 org.reestr.r709.total.totalSng += org.SngCallsTotal.Man
 org.reestr.r709.total.totalTkm += org.TkmCallsTotal.Man
 org.reestr.r709.total.totalLoc += org.LocalCallsTotal.Man
 org.reestr.r709.total.totalPstn += org.PstnServiceTotal.Man
 org.reestr.r709.total.totalTen += org.tenPercentMan
 org.reestr.r709.total.totalMan += org.orgTotalMan
 fileMutex.Lock()
 defer fileMutex.Unlock()
 org.reestr.r709.file.WriteString(str)
 return
 fmt.Println("REESTR",org.IntCallsTotal) //TotalDual
 fmt.Println("REESTR",org.SngCallsTotal) //TotalDual
 fmt.Println("REESTR",org.TkmCallsTotal) //TotalDual
 fmt.Println("REESTR",org.LocalCallsTotal) //TotalDual
 fmt.Println("REESTR",org.PstnServiceTotal) //TotalDual

 fmt.Println("REESTR",org.tenPercentUsd) //float64
 fmt.Println("REESTR",org.tenPercentMan) //float64

 fmt.Println("REESTR",org.orgTotalUsd) //float64
 fmt.Println("REESTR",org.orgTotalMan) //float64

}

func alignCenter(s string,w int) string {
 runes := []rune(s)
 length := len(runes)
 return fmt.Sprintf("!%[1]*s", -w, fmt.Sprintf("%[1]*s", (w + length)/2, s))
}

func alignLeft(s string,w int) string {
 runes := []rune(s)
 length := len(runes)
 str := string(runes[0:w-1])
 fmt.Println("alignLeft",s,w,length,str)
 return fmt.Sprintf("! %[1]*s", -w, str)
}

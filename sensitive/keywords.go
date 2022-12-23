package sensitive

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var sendatapath string
var sentivedir = "sensitive"
var sensitiveTxt = "sensitive.txt"

func SernsitiveInit(dataPath string) {
	//SenSitiveStatus = 1
	//mkdriDht(dataPath)
	
	personSensitiveWords()
	go initSystemSensitiveWords()
}

//sensitive,
func mkdriDht(dataPath string) {
	sendatapath = dataPath
	senTextPath := filepath.Join(dataPath, sentivedir)
	_, err := os.Stat(senTextPath)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(senTextPath, os.ModePerm)
		if err != nil {
			logrus.Error("", err)
		}
	}
}


func SensitiveSave(sensitiveWords string) error {
	return writeSensitiveToFile(sensitiveWords)
}


func QuerySensitive() (string, error) {
	sentiByte, err := readSensitiveWords()
	if err != nil {
		logrus.Error("", err)
	}
	if sentiByte != nil {
		return string(sentiByte), err
	} else {
		return "", nil
	}

}


func writeSensitiveToFile(sensitiveWords string) error {
	btTextPath := filepath.Join(sendatapath, sentivedir, sensitiveTxt)
	f, err := os.OpenFile(btTextPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		if !strings.Contains(err.Error(), "The system cannot find the file specified") && !strings.Contains(err.Error(), "no such file or directory") {
			logrus.Error("", err)
			return err
		} else {
			f, err = os.Create(btTextPath)
			if err != nil {
				logrus.Error("", err)
				return err
			}
		}
	}
	defer f.Close()
	n, err := f.Write([]byte(sensitiveWords))
	if err == nil && n < len([]byte(sensitiveWords)) {
		err = io.ErrShortWrite
	}
	if err != nil {
		logrus.Error("", err)
		return err
	}
	return nil
}


func readSensitiveWords() ([]byte, error) {
	//filePth := filepath.Join(sendatapath, sentivedir, sensitiveTxt)
	//f, err := os.Open(filePth)
	//if err != nil {
	//	if strings.Contains(err.Error(), "The system cannot find the file specified") || strings.Contains(err.Error(), "no such file or directory") {
	//		
	//		sensitiveWords := strings.Replace(personDefaultWord, " ", SensitiveWordSplit, -1)
	//		writeSensitiveToFile(sensitiveWords)
	return []byte(personDefaultWord), nil
	//	} else {
	//		return nil, err
	//	}
	//}
	//return ioutil.ReadAll(f)
}

const systemword = "                        JU HUANG HUANGJU huang ju huangju                                                                                                          ·                                                                                                                                                                       ·                                                           ·              ·                                                                                                                     ·                                                                                                                                                                                                                                                                                                             Hui      310                                        JIANZEMIN LIPENG FALUN  LIHONGZH jiangzemin B    PowertotheFalunGong BrothersinArms theUndergroundResistance FalunGong UndergroundResistanceArms BrothersTIBET tibet              LIHONGZHI                1 XJG351                                                    ·         2.23                     001                        5.4     《：》    **                     CHENG                                                                                                                                                                                                                                                                                                                                                       16  64  64 89                                a-lun-gong FLG                                                                                core    ?                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  89 89 9ping 9 9 9 fa fa FL freetibet GONG GONG HEIDEROOSJESAND LOTENUNITEDTEBET hyperballad-tibet jiuping UnitedTibet      0 lun X     china cn     ping                                            "

const personDefaultWord = "1|-|                                        B                       91 BJ Blow Job CAR SEX G  KJ Y     after-play          yanjiaoshequ                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      BB                                                                                        ~                                                                                                                                                                                                                                                                                                                                 □                                                                                         av japansweet adult 16dy-             999               anal        incest       B         97sese AV  33bbb   BlowJobs      cc  cc                     a4u         99bb g tw18 asiasex teen sexy       A                 hardcore amateur            hotsex porn                 uu   mm  77bbb      A p     97         3p   AV                AV  B           a4y     P    18 g                  H        petgirl        55sss  xiao77      222se    xx  bt    の           xx                                                 [hz] (hz) [av] (av) [sm] (sm) sm      sb     AIDS aids Aids DICK dick Dick penis sex SM                Bitch cao FUCK Fuck kao NMD NND SHIT SUCK Suck tnnd K                                  B                         PENIS BITCH BLOWJOB KISSMYAS    X                                            b                 H hgame H   shit B b                     A A                                                                                                   asia_fox        3P    JB  _ X           "

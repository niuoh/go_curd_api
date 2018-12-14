package function

import(
  //  "strings"
    "unicode"
    "time"
    "fmt"
    "math/rand"
)

func Ucfirst(str string) string {
    for i, v := range str {
        return string(unicode.ToUpper(v)) + str[i+1:]
    }
    return ""
}

func InArray(need string, needArr []string) bool {
     for _,v := range needArr{
        if need == v{
            return true
        }
    }
    return false
}

func HaveKey(obj map[interface{}]interface{},key string) bool{
  for k,_ := range obj{
    if(k==key){
      return true
    }
  }
  return false
}

func RandString(len int) string {
    r := rand.New(rand.NewSource(time.Now().Unix()))
    bytes := make([]byte, len)
    var b int
    for i := 0; i < len; i++ {

      m:=r.Intn(3)
      if(m==0){
        b = r.Intn(27) + 97
      }else if(m==1){
        b = r.Intn(27) + 65
      }else if(m==2){
        b = r.Intn(10) + 48
      }
        bytes[i] = byte(b)
    }
    return string(bytes)
}

func GetTime() string {
    now := time.Now()
    year, mon, day := now.Date()
    hour, min, sec := now.Clock()
    return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", year, mon, day, hour, min, sec)
}

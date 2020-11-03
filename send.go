package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://www.zhihu.com/hot"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", " _zap=471de525-fba6-41c5-b5da-ec8314e31480; d_c0=\"ALCd5VTyvxGPTsw7M9ZyFIDSYT2Z_IOxge8=|1597719094\"; Hm_lvt_98beee57fd2ef70ccdd5ca52b9740c49=1604404536,1604405029,1604406162,1604407385; _ga=GA1.2.795880107.1597719097; capsion_ticket=\"2|1:0|10:1604389506|14:capsion_ticket|44:YTMyMzY0NzhiZGYxNDRkYWI2ZmEzNGRmMTEyNmQ1YTk=|802684f7e3f8c977f2c5a72fcbd8a2cb5d56c6082efab195011968b3f6ccc9bc\"; _xsrf=hcBsx6VvOV6yoYnzT621R2DJwYCbwgOh; Hm_lpvt_98beee57fd2ef70ccdd5ca52b9740c49=1604407385; l_n_c=1; r_cap_id=\"NmNkODI5MTJmMzA2NGQwNDhlODY4NTZkNjI0YzgzNTU=|1604389512|d1bccf105a4ca4ff75b758258a1e319374582a19\"; cap_id=\"OTYzMmMzMGY2ZmNmNDQ4YmEwYjU2YWUzOGE1ZTcwYmU=|1604389512|0067be4cd11e81292235c987bdf0717d596a83a1\"; l_cap_id=\"ZmQyY2ZkMDY1NTlhNGM3ZDk3OGM2Mzg3YTdiY2IwMjI=|1604389512|dc33462ecb45ff74519f981381fa16b3f3c39730\"; n_c=1; z_c0=Mi4xbEZvYUR3QUFBQUFBc0ozbFZQS19FUmNBQUFCaEFsVk5qbGlPWUFDY3VDQXl5S0NMUVVZczRTV0JZcWk1eDV3cFR3|1604389518|95ab58b527f650fd47f0ff6feb9295ef7cddedda; tst=h; KLBRSID=fb3eda1aa35a9ed9f88f346a7a3ebe83|1604406163|1604404020; tshl=; q_c1=7ae97f24b5844158a63a508677648097|1604392022000|1604392022000; _zap=eb58e67e-2071-41a5-9eb8-aa0731680631; _xsrf=huHWSQwxScw2VxvC60wxcXw2tYBRSmFt; KLBRSID=fb3eda1aa35a9ed9f88f346a7a3ebe83|1604407472|1604404020")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

package conv

import (
	"bytes"
	"crypto/md5"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestUnicode2Utf8(t *testing.T) {
	testCases := map[uint32]string{
		0x20:    " ",
		0x67:    "g",
		0xA9:    "©",
		0x8427:  "萧",
		0x1F604: "😄",
	}

	for k, v := range testCases {
		if u8, err := unicodeToUtf8(k); err != nil {
			t.Fatal(err)
		} else {
			if string(u8) != v {
				t.Fatalf("unpass: %x, got: %x, except: %x", k, u8, []byte(v))
			}
		}
	}
}

func TestGbk2Utf8(t *testing.T) {
	testCases := map[uint16]string{
		0xcff4: "萧",
		0xd5c5: "张",
		0xe9cf: "橄",
		0xa3fe: "￣",
	}

	for k, v := range testCases {
		if u8, err := gbkToUtf8(k); err != nil {
			t.Fatal(err)
		} else {
			if string(u8) != v {
				t.Fatalf("unpass: %x, got: %x, except: %x", k, u8, []byte(v))
			}
		}
	}
}

func TestGbk2Utf8Multi(t *testing.T) {
	var (
		src       string = "data/936.txt"
		md5CmdOut bytes.Buffer
	)

	iconvCmd := exec.Command("iconv", "-c", "-f", "gbk", "-t", "utf-8", src)
	md5Cmd := exec.Command("md5")
	md5Cmd.Stdin, _ = iconvCmd.StdoutPipe()
	md5Cmd.Stdout = &md5CmdOut

	md5Cmd.Start()
	iconvCmd.Run()
	md5Cmd.Wait()

	if file, err := os.Open(src); err != nil {
		t.Fatal(err)
	} else {
		defer file.Close()

		var out bytes.Buffer
		GbkToUtf8(file, &out, true)

		ioutil.WriteFile("data/conv_test.txt", out.Bytes(), 0666)

		md5Crp := md5.New()
		md5Crp.Write(out.Bytes())
		md5CrpSum := string(md5Crp.Sum(nil))

		if md5CmdOut.String() != md5CrpSum {
			t.Fatalf("md5 checksum aren't equal %x", md5Crp.Sum(nil))
		}
	}
}

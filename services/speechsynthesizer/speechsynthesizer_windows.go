package speechsynthesizer

import "bytes"
import "text/template"
import "os/exec"
import "os"

const script = `
Add-Type -AssemblyName System.speech
$n = New-Object System.Speech.Synthesis.SpeechSynthesizer

$n.Rate = {{.Rate}}

{{with .Voice}}
$n.SelectVoice("{{.}}")
{{end}}

$n.Speak("{{.Text}}")
`

type Notification struct {
	Text string
	// Rate is from -10 to 10. -10 is slowest.
	Rate  int
	Voice string
}

func (n *Notification) Send() error {
	tmpl, err := template.New("").Parse(script)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, n); err != nil {
		return err
	}

	cmd := exec.Command("PowerShell", "-Command", buf.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}


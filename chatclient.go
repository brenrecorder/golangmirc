package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

"net/http"
"io/ioutil"

"time"
"runtime"
"os/exec"
"bufio"
)

var serveraddress string = "127.0.0.1"
var nickname string

var servername string
var mode = "exit"
var menuchoice int = 0
var lengthchatmsg int = 25
const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("12"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }




func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	//str := fmt.Sprintf("%d. %s", index+1, i)
str := fmt.Sprintf("%s", i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("●  " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			os.Exit(0)
		case "ctrl+c":
			m.quitting = true
			os.Exit(0)
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {

if (m.choice == "Join server") { menuchoice = 1 } 
if (m.choice == "Start server") { menuchoice = 2 } 
if (m.choice == "Settings") { menuchoice = 3 } 
if (m.choice == "Help") { menuchoice = 4 } 
if (m.choice == "Exit") { menuchoice = 5 } 

	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Menu: %s", m.choice))
	}
	if m.quitting {
		os.Exit(0)
		return quitTextStyle.Render("Quitting program")
	}
	return "\n" + m.list.View()
}

type (
	errMsg error
)


func main() {

getsettings()
menuchoice = 0
servername=""
if (len(nickname) < 1) {
fmt.Print("Chat\n\nYour nickname: ")
fmt.Scanln(&nickname)
setsettings()
if (len(nickname) < 2) { 
fmt.Print("\nPlease enter a nickname of 2 or more characters.\n")
nickname = "" 
main()
 }
}

ClearTerminal()

	items := []list.Item{
		item("Join server"),
		item("Start server"),
		item("Settings"),
		item("Help"),
		item("Exit"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "SpecChat\nNickname: "+nickname+"\n\nMenu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
if (menuchoice != 1 && menuchoice != 2 && menuchoice != 3 && menuchoice != 4 && menuchoice != 5) { main() }
if (menuchoice == 1) {
ClearTerminal()
fmt.Print("Chat server join\n\nServer name\t\t\t\tMessages\n")
   resp, err := http.Get("http://" + serveraddress + "/chatserver.php?action=serverlist")
   if err != nil {
      fmt.Println(err)
	  main()
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Println(err)
   }
   sbody := string(body)
   serverlist := strings.Split(sbody, "-") 
      if (len(serverlist) > 0) {
      for _, server := range serverlist {
	    servername := strings.Split(server, ":") 
		servername[0] = strings.ReplaceAll(servername[0], "_", " ")
		if (len(servername) == 2) {
		if (len(servername[0]) < 25) {
		for (len(servername[0]) < 25) { servername[0] = servername[0] + " " }
		}
		fmt.Print("\n" + servername[0] + "\t\t" + servername[1])
		}
	  }
	  
	  }
	  servername=""
   fmt.Print("\n\nJoin server: ")
scanner := bufio.NewScanner(os.Stdin)
if scanner.Scan() {
    servername = scanner.Text()
}
if (len(servername) < 1) { main() }
servernameb := strings.Replace(servername, " ", "_", -1)
mode = "chat"
serversync(servernameb)
 }
 //echo $row['ServerName'].":".$row['Messages']."-";
if (menuchoice == 2) {
ClearTerminal()
servername=""
fmt.Print("Chat server start\n\nServer name: ")
scanner := bufio.NewScanner(os.Stdin)
if scanner.Scan() {
    servername = scanner.Text()
}
if (len(servername) < 1) { main() }
var pwdserver string
fmt.Print("Password: ")
fmt.Scanln(&pwdserver)
pwdserver = strings.Replace(pwdserver, "\n", "", -1)
servernameb := strings.Replace(servername, " ", "_", -1)
_,err := http.Get("http://" + serveraddress + "/chatserver.php" + "?action=createserver&ServerName="+servernameb + "&Nickname=" + nickname + "&Password=" +pwdserver)
   if err != nil {
      fmt.Println(err)
	  main()
   }
mode = "chat"
serversync(servernameb)
}
if (menuchoice == 3) { 
fmt.Print("Chat Settings\n\nServer name: " + serveraddress)
setsettings()
main()
 }
 if (menuchoice == 4) { 
ClearTerminal()
fmt.Print("Chat help\n\nChoose 1 to join server\nChoose 2 to start server\n\nEnter 'quithost' to stop your chatserver\nEnter 'exit' to join another server\n\n") 
fmt.Scanln()
main()
}
if (menuchoice == 5) { os.Exit(0) }
}




func getsettings() {
readsettings, err := os.ReadFile("settings.ini")
if err != nil {
fmt.Print(err)
}
strSettings := string(readsettings)
stringSettings := strings.ReplaceAll(strSettings, "\n", "")

splitsettings := strings.Split(stringSettings, ":")
if (len(splitsettings) > 1) {
serveraddress = splitsettings[0]
nickname = splitsettings[1]
} else {
serveraddress = "127.0.0.1"
}
return
}
func setsettings() {
var newserver string
var newnickname string
fmt.Print("\nMain server address: ")
fmt.Scanln(&newserver)
fmt.Print("\nNickname ("+nickname+"): ")
fmt.Scanln(&newnickname)
if (len(newnickname) >2) {
nickname = newnickname
}
if (len(newserver) >2) {
val := newserver
data := []byte(val + ":" + nickname)
err := ioutil.WriteFile("settings.ini", data, 0)
if err != nil {
    fmt.Println(err)
}
serveraddress = newserver
fmt.Println("done setting main server " + serveraddress)
} else {
fmt.Println("set old main server " + serveraddress)
}
return
}


var chatsold int = 0

func serversync(nameserver string) {
nameserverb := strings.Replace(nameserver, " ", "_", -1)
for (mode == "chat"){
go receivechat(nameserverb)
time.Sleep(1 * time.Second)
}
servername = ""

 main()
}
func sendnewchatmsg(nameserver string, newchatmessage string) {
nameserverb := strings.Replace(nameserver, " ", "_", -1)

sendchat := http.Client{
    Timeout: 5 * time.Second,
}

_,err := sendchat.Get("http://" + serveraddress + "/chatserver.php" + "?action=sendmsg&ServerName="+nameserverb + "&Nickname=" + nickname + "&Message=" + newchatmessage)
   if err != nil {
      fmt.Println(err)
   }
return
}
func receivechat(nameserver string) {
nameserverb := strings.Replace(nameserver, " ", "_", -1)
nameserverconn := strings.Replace(nameserver, "_", " ", -1)
//127.0.0.1/chatserver.php?action=sendmsg&ServerName=ChatServerA&Nickname=Spec12&Message=hoi
for (mode == "chat"){
   resp, err := http.Get("http://" + serveraddress + "/chatserver.php?action=viewmsg&ServerName="+nameserverb)
   if err != nil {
      fmt.Println(err)
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Println(err)
   }
   sbody := string(body)
   messages := strings.Split(sbody, "-") 
   if (len(messages) != chatsold) {
   ClearTerminal()
   fmt.Print("Connected to: " + nameserverconn + "\n")
   var cntermax int =50
   var cnter int =0
      for _, message := range messages {
	  if (cnter > (len(messages)-cntermax)) {
	    messagepart := strings.Split(message, ":") 
		if (len(messagepart) == 2 && len(messagepart[1])>0) {
		var newnicksized string = messagepart[0]
		for (len(newnicksized) < 15) { newnicksized = newnicksized + " " }
		var cntletters int = 0
		var cntlettersb int = 0
		var messagesized string 
		var currmessage string = messagepart[1]
		if (len(currmessage) > lengthchatmsg) {
		for (cntletters < len(currmessage)) {
		messagesized = messagesized + string(currmessage[cntletters])
		if (cntlettersb > lengthchatmsg) { 
		messagesized = messagesized + "\n\t\t " 
		cntlettersb = 0 }
		cntlettersb++
		cntletters++
		}
		} else {
		messagesized = messagesized + string(currmessage)
		}
		fmt.Print("\n" + newnicksized + ": " + messagesized)
		}
		} else { cnter++ }
	  }
	  fmt.Print("\n\n")
	  var newmessage string = ""
	  var scan string
	  chatsold = len(messages)
	  fmt.Print("Chat: ")
	
		scanner := bufio.NewScanner(os.Stdin)
scanner.Scan()
scan = scanner.Text()

		newmessage = strings.Replace(scan, "'", " ", -1)
	
	  if (newmessage == "quithost") { 
		  fmt.Print("\nEnter password: ")
		  var quithostpassword string
		  fmt.Scanln(&quithostpassword)
		  _,err := http.Get("http://" + serveraddress + "/chatserver.php" + "?action=deleteserver&ServerName="+nameserverb + "&Password=" + quithostpassword)
			if err != nil {
				fmt.Println(err)
			} else {
			  fmt.Print("\nChat host stop request sent..")
			   mode = "exit"
			  break
			}
	  }
	  if (newmessage == "exit") {  
	   mode = "exit"
	   } else {
	  if (len(newmessage) > 0) {sendnewchatmsg(nameserver, newmessage) } }
	  
	  break
	  } else {
	  break
	  }
	  //fmt.Print("\nChat:	")
	  //chatmsg := fmt.Scanln()
	 
}
return
//echo $row['Nickname'].":".$row['Message']."-";
}

func runCmd(name string, arg ...string) {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Run()
}
func ClearTerminal() {
    switch runtime.GOOS {
    case "darwin":
        runCmd("clear")
    case "linux":
        runCmd("clear")
    case "windows":
        runCmd("cmd", "/c", "cls")
    default:
        runCmd("clear")
    }
}
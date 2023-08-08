<?php
$sqlNewTableChats = "CREATE TABLE ChatServers (
ID INTEGER PRIMARY KEY AUTOINCREMENT,
ServerName TEXT,
Messages INTEGER,
Nickname TEXT,
Message TEXT,
Password TEXT,
DATE TEXT
)";
$pattern = '/[^a-zA-Z0-9@\(.*?\)@%^&!\']/';

if (file_exists("chatserver1.db")) { //CREATE DATABASE
	$db = new SQLite3('chatserver1.db');
} else {
	$db = new SQLite3('chatserver1.db');
	$db->querySingle($sqlNewTableChats);


	echo "Tables and file chats created..";
}

if (!empty($_GET['action'])) { $action = preg_replace($pattern, '', $_GET['action']); } else { $action = ""; }


if ($action == "createserver") {;
	if (!empty($_GET['ServerName'])) { $ServerName = preg_replace($pattern, '', $_GET['ServerName']); } else { $ServerName = ""; }
	if (!empty($_GET['Nickname'])) { $Nickname = preg_replace($pattern, '', $_GET['Nickname']); } else { $Nickname = ""; }
	if (!empty($_GET['Password'])) { $Password = md5($_GET['Password']); } else { $Password = ""; }
	if (!empty($_GET['delete'])) { $delete = preg_replace($pattern, '', $_GET['delete']); } else { $delete = ""; }
	
	if (empty($delete)) {
	$db->exec("INSERT INTO ChatServers(ServerName, Messages, Nickname, Message, Password, DATE) VALUES('".$ServerName."', '0', '".$Nickname."', '', '".$Password."', '".date("Y-m-d")."')");	
	} else {
	$db->exec("DELETE FROM ChatServers WHERE ServerName=='".$ServerName."'");	
	}
}
if ($action == "deleteserver") {
	if (!empty($_GET['ServerName'])) { $ServerName = preg_replace($pattern, '', $_GET['ServerName']);  } else { $ServerName = ""; }
	if (!empty($_GET['Password'])) { $Password =  preg_replace($pattern, '', $_GET['Password']);  } else { $Password = ""; }
	$pwdmd5 = md5($Password);
	echo $pwdmd5;
	$count = $db->querySingle("SELECT COUNT(ServerName) as count FROM ChatServers WHERE Password ='".$pwdmd5."'");
	echo $count;
	if ($count > 0) {
		$db->exec("DELETE FROM ChatServers WHERE ServerName == '".$ServerName."'");	
		echo "server deleted : " . $ServerName;
	}
}
if ($action == "sendmsg") {

	if (!empty($_GET['ServerName'])) { $ServerNameB = preg_replace($pattern, '', $_GET['ServerName']);  } else { $ServerNameB = ""; }
	
	if (!empty($_GET['Nickname'])) { $Nickname = preg_replace($pattern, '', $_GET['Nickname']);  } else { $Nickname = ""; }
	if (!empty($_GET['Message'])) { $Message = preg_replace($pattern, '', $_GET['Message']); } else { $Message = ""; }


	
	$cnt =0;
	$db->exec("INSERT INTO ChatServers(ServerName, Messages, Nickname, Message, DATE) VALUES('".$ServerNameB."', '".$cnt."', '".$Nickname."', '".$Message."', '".date("Y-m-d")."')");
	$cntMSG = $db->querySingle("SELECT Messages as cnter FROM ChatServers WHERE ServerName=='".$ServerNameB."'");
	$db->exec("UPDATE ChatServers SET Messages='".($cntMSG+1)."' WHERE ServerName='".$ServerNameB."'"); 
	


//}

	//if($cntMSG >5) {
	//	$db->exec("DELETE FROM ChatServers WHERE ROWID IN (SELECT ROWID FROM ChatServers ORDER ROWID DESC LIMIT -1 OFFSET 5)");	
	//}
}  
if ($action == "cleanup") {
$count = $db->querySingle("SELECT COUNT(ServerName) as count FROM ChatServers");
if ($count > 1) {
$idin = $db->querySingle("SELECT min(id) as test FROM ChatServers");
$db->exec("DELETE FROM ChatServers WHERE ID == '".$idin."'");	
echo $idin;
}
}
if ($action == "viewmsg") {
if (!empty($_GET['ServerName'])) { $ServerName = preg_replace($pattern, '', $_GET['ServerName']); } else { $ServerName = ""; }
$query="SELECT ServerName, Nickname, Message FROM ChatServers WHERE ServerName = '".$ServerName."'";
$result=$db->query($query);
while($row= $result->fetchArray()){
echo $row['Nickname'].":".$row['Message']."-";
}

}

if ($action == "serverlist") {

$query="SELECT DISTINCT ServerName, Messages FROM ChatServers";
$result=$db->query($query);
while($row= $result->fetchArray()){
if (strlen($row['ServerName']) > 0) {
echo $row['ServerName'].":".$row['Messages']."-"; }
}

}
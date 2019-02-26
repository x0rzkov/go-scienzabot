package consts

//InitSQLString is the initialization query, to run once at the bot startup
const InitSQLString = `
/*
The Users table is supposed to contain all the users subscribed to the bot
*/
CREATE TABLE IF NOT EXISTS 'Users' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Nickname'  TEXT UNIQUE,
	'Biography'  TEXT,
	'Status'  INTEGER NOT NULL DEFAULT 0,
	'Locale'  TEXT NOT NULL DEFAULT '` + DefaultLocale + `',
	'Permissions'  INTEGER NOT NULL DEFAULT 0,
	'LastSeen'  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	'RegisterDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

/*
The Groups table is supposed to contain all the groups the bot is added to.
Status refers if the bot was kicked form the group, not used atm
*/
CREATE TABLE IF NOT EXISTS 'Groups' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Title'  TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,
	'Locale'	TEXT NOT NULL DEFAULT '` + DefaultLocale + `',
	'Status'	INTEGER NOT NULL DEFAULT 0
);

/*
The Channels table is supposed to contain the channels that admins may want to use to forward messages
from a group to a channel referring to a particular message
*/
CREATE TABLE IF NOT EXISTS 'Channels' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'GroupID'  INTEGER NOT NULL,
	'Name'	TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,	
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_channels_channel_group__unique UNIQUE ('ID','GroupID')
);

/*
The Permissions table is supposed to contain the permissions for each user in each group.
*/
CREATE TABLE IF NOT EXISTS 'Permissions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'	INTEGER NOT NULL,
	'GroupID'	INTEGER NOT NULL,
	'Permission' INTEGER DEFAULT 0,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_perm_user_group_perm_unique UNIQUE ('UserID','GroupID','Permission')
);

/*
The Lists table is suopposed to contain the lists where a user can subscribe to.
Such list should be group-dependent (if not specified otherwise on the status field, that shouold be based on a bit-based flag)
The status is not used yet
*/
CREATE TABLE IF NOT EXISTS 'Lists' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'Name'  TEXT NOT NULL,
	'GroupID'	INTEGER NOT NULL,
	'Properties'  INTEGER DEFAULT 0,
	'CreationDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	'LatestInvocationDate' TIMESTAMP,
	'Parent' INTEGER,
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	FOREIGN KEY('Parent') REFERENCES Lists('ID'),
	CONSTRAINT con_lists_name_group_unique UNIQUE ('Name','GroupID')
);

/*
The Bookmarks table is used to when a user wants to save a message for a future reference.
The bot will in fact save the group and the message, and will bind it to a user.
The bot will also save a copy of the message content (when possible).
Deletion of a row should be impossibilitated to a user
TODO: Remembere to check if the message still exists
*/
CREATE TABLE IF NOT EXISTS 'Bookmarks' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'  INTEGER NOT NULL,
	'GroupID'  INTEGER NOT NULL,
	'MessageID'	INTEGER NOT NULL,
	'Alias'	TEXT,
	'Status' INTEGER DEFAULT 0,
	'MessageContent' TEXT, 
	'CreationDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	'LastAccessDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_bookm_user_group_unique UNIQUE ('UserID','GroupID')

);

/*
The Subscriptions table is used to subscribe a specific user to a "list" where he belongs
*/
CREATE TABLE IF NOT EXISTS 'Subscriptions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'ListID'  INTEGER NOT NULL,
	'UserID'  INTEGER NOT NULL,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('ListID') REFERENCES 'Lists'('ID'),
	CONSTRAINT con_subs_user_list_unique UNIQUE ('UserID','ListID')
);

/*
The MessageCount table is used to count the message of each user in the various groups
This allows the bot to count the message of a specific user on a multitude of groups
*/
CREATE TABLE IF NOT EXISTS 'MessageCount' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'  INTEGER NOT NULL,
	'GroupID'  INTEGER NOT NULL,
	'MessageCount'  INTEGER NOT NULL,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_msgcoubt_user_group_unique UNIQUE ('UserID','GroupID')
);

/*
The Strings table will contain all the strings about a specific group 
(in fact it's group-dependent).
Such strings could be something like a welcome message, a help message an so on...
*/
CREATE TABLE IF NOT EXISTS 'Strings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL,
	'Value'	TEXT DEFAULT 'Not implemented',
	'Locale'	TEXT DEFAULT '` + DefaultLocale + `',
	'GroupID'	INTEGER NOT NULL,
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_strings_key_group_locale_unique UNIQUE ('Key','GroupID','Locale')
);

/*
The Settings table will contain all the settings about a specific group 
(in fact it's group-dependent).
An example could be the status (on/off) of a specific function of the bot
*/
CREATE TABLE IF NOT EXISTS 'Settings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL,
	'Value'	TEXT DEFAULT 'Not implemented',
	'GroupID'	INTEGER NOT NULL,
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_setting_key_group_unique UNIQUE ('Key','GroupID')
);

/*
The BotSettings table will contain all the settings of the bot
Such settings could be things like the default locale
TODO: evaluate UNIQUE on Key
TODO: evaluate removal of this table
*/
CREATE TABLE IF NOT EXISTS 'BotSettings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL UNIQUE,
	'Value'	TEXT DEFAULT 'Not implemented'
);

/*
The BotStrings table will contain all the strings to be used from the bot, like the "cancel" text and so on...
As the constraint shows, there can only bae a pair of key-locale per table(we can't have 2 way of saying the same thing
in the same language; which one should we take?) 
*/
CREATE TABLE IF NOT EXISTS 'BotStrings' (
	'ID'		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'		TEXT NOT NULL,
	'Value'		TEXT DEFAULT 'Not implemented',
	'Locale'	TEXT DEFAULT '` + DefaultLocale + `',
	CONSTRAINT con_botstrings_key_locale_unique UNIQUE ('Key','Locale')
);

/*
The Log table should contains the log of the bot, like a new user subscribed, a help message requested,
a non-working command received and so on.
Also crashes should be documented.
The Event field is supposed to be a human-readable value
The RelatedUser field is supposed to contain the user in the context who sent or triggered the action
The RelatedGroup field should contain the group where the action triggered
The Message field should be a human-readble string that should say what happened (moreless)
The UpdateValue field should contain the update raw string from telegram
The Error field should contain the string reported by the funcion that throwed the error (if any)
The Severity field should indicate he gravity of the event
The Date field is the date when the event occourred
*/
CREATE TABLE IF NOT EXISTS 'Log' (
	'ID'		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Event'		TEXT NOT NULL,
	'RelatedUser'		INTEGER,
	'RelatedGroup'		INTEGER,
	'Message'		TEXT,
	'UpdateValue'		TEXT,
	'Error'		TEXT,
	'Severity'		INTEGER,
	'Date'	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inserting the default locale in DB
INSERT OR IGNORE INTO BotSettings (Key, Value ) VALUES ( "DefaultLocale", "'` + DefaultLocale + `'" );

-- Inserting Pandry and AndreaIdini as users
INSERT OR IGNORE INTO Users (ID, Nickname, Permissions) VALUES (14092073, "Pandry", 255), (44917659, "AndreaIdini", 255);

-- Inserting bot version if not exists
INSERT OR IGNORE INTO BotSettings (Key,Value) VALUES ("version", "⚛️ v 0.1g α");

-- Message to ask used to use the command in private chat
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("onPrivateChatCommand","Please, ask me that in private chat", "en");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("onPrivateChatCommand","Perfavore, usa questo comando in chat privata", "it");

-- Help message
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("helpCommand","Comandi per aiuto:
/help - Don't ya know?", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("helpCommand","Help commands:
/help - Don't ya know?", "en");

-- Info Command
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("infoCommand","Ciao 😁
Sei confuso? 
Questo è il bot del gruppo @scienza e permette di usufruire di queste funzioni:
/iscriviti per iscriverti al database di utenti e per partecipare alle liste sugli interessi
/aderisci per iscriverti ad una lista, puoi usare anche: /partecipa, /registrati e /sottoscrivi
/bio per scrivere qualcosa su di te
/liste per scoprire le liste già presenti
/gdpr consulta le norme sul GDPR
/privs elenca i privilegi di un utente (richiesto come argomento)
/biografia mostra la biografia di un utente (richiesto come argomento)
/disiscrivi per cancellarti da una lista alla quale hai aderito, puoi usare pure: /esci, /rimuovi, /iscrizioni e /aderenze
/info Ottieni informazioni su di me
Puoi anche usare il bot in modalità ""inline"": sarà sufficiente scrivere @scienziati_bot <username> per avere informazioni riguardo l'utente
In caso di problemi invece, sei pregato di conttattare @Pandry, in quanto sviluppatore del bot.
Report di problemi, come ad esempio liste non presenti, bot non responsivo ecc sono assolutamente gradite; O anche solo per proporre qualche idea e conversarne a riguardo.
A tal proposito, esiste un gruppo dedicato ai programmi scritti in comune tra i membri di @Scienza.
Chiedi ad un amministratore per ulteriori informazioni a riguardo.
Buona continuazione su @Scienza", "it");

INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("infoCommand","That's the @scienza custom bot and things", "en");

-- Errors strings
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("generalError","❌ Error 518 ❌ - Bip bop, I'm a teapot
Si è verificato un errore.
Lo sviluppatore (@Pandry) è stato già avvertito.
Sei pregato di contattarlo per descrivere in che modo questo errore è stato visualizzat, grazie.", "it");

INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("generalError"," ❌ Error 518 ❌ - Bip bop, I'm a teapot
An error occourred.
The developer (@Pandry) has been notified.
You are kindly asked to text him telling what you've done to see this.", "en");


-- User added strings
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("userAddedSuccessfully","✅ Complimenti, ti sei iscritto con successo!
Ora puoi fare uso delle funzionalità del bot, buona continuazione!", "it");

INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("userAddedSuccessfully","✅ Congratulations, you registred successfully!
You can now use all the features of the bot! Enjoy you time :3", "en");

-- User already registred
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("userAlreadyRegistred","⚠️ Attenzione!
Risulti già iscritto al bot, non è necessario tentare nuovamente di iscriversi.", "it");

INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("userAlreadyRegistred","⚠️ Warning!
It looks like you are already registred to the bot, you don't need to register again.", "en");

-- User already registred
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("userNotRegistred","⚠️ Attenzione!
Per fare uso di questa funzionalità è necessario iscriversi.", "it");

INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("userNotRegistred","⚠️ Warning!
To use this feature you need to register.", "en");

-- Delete message button text
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("deleteMessageText","🗑 Elimina questo messaggio 🗑", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("deleteMessageText","🗑 Delete this message 🗑", "en");

-- Close message button text
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("closeMessageText","🗑 Chiudi", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("closeMessageText","🗑 Close", "en");

-- List created successfully
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("listCreatedSuccessfully","✅  Lista creata con successo", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("listCreatedSuccessfully","✅  List created successfully", "en");

-- Newlist Syntax error
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("newlistSyntaxError"," ❗️ Errore di sintassi - l'uso previsto è /newlist <nomelista>
Il nome deve contenere solo caratteri minuscoli, trattini e underscores [a-z\-_]{1,30}, senza spazi e fino a 30 caratteri.", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("newlistSyntaxError"," ❗️ Syntax error - the use is supposed to be /newlist <listname>
The name shall only contains lowercase characters, dashed and underscores [a-z\-_]{1,30}, without spaces and up to 30 chars.", "en");

-- notAuthorized
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("notAuthorized"," 🚷 Errore 403 - permessi insufficienti,
Se si ritiene che questo sia un errore, si è gentilmente pregati di contattare l'amministratore del bot (@Pandry)", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("notAuthorized"," 🚷 Error 403 - Unauthorized,
If you believe this is an error, please contact the bot author (@Pandry)", "en");


--notImplemented
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("notImplemented"," Questa funzionalità non è ancora implementata", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("notImplemented"," This feature is not implemented yet", "en");


-- Available Lists
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("availableLists","I pulsanti sottostanti rappresentano le liste disponibili.
Tappa su una lista per iscriverti", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("availableLists","The ones below are the availalbe lists.
Tap on one of them to subscribe to them.", "en");

-- callbackQueryAnswerSuccess
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("callbackQueryAnswerSuccess","✅ Successo!", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("callbackQueryAnswerSuccess","✅ Success!", "en");

-- callbackQueryAnswerError
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("callbackQueryAnswerError","❌ Errore
Perfavore, contattare @Pandry", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("callbackQueryAnswerError","❌ Error
Please contact @Pandry", "en");

-- callbackQueryAnswerError
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("noListsLeft","❌ Errore
Non sono presenti liste che a cui puoi iscriverti", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("noListsLeft","❌ Error
It aint no list you can subscribe to", "en");

-- noSubscription
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("noSubscription","❌ Errore
Nonn ti sei ancora iscritto a nessuna lista", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("noSubscription","❌ Error
YOu did not join any list yet", "en");

-- subscribedLists
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("subscribedLists","Di seguito puoi trovare le liste a cui sei iscritto.
Per revocare una sottoscrizione, sei pregato di ""tappare"" la lista in questione.", "it");
INSERT OR IGNORE INTO BotStrings (Key, Value, Locale) VALUES ("subscribedLists","Here you can see the lists you are currently subscribed to.
To remove a subscription you can ""tap"" on the list.", "en");

`

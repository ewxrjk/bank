package main

// Automatically generated code, don't edit
const appcss = ("body {\n" +
	"    font-family: Liberation;\n" +
	"    font-size: 12pt;\n" +
	"}\n" +
	"div.transactions { width: max-content }\n" +
	"table.transactions { border-collapse: collapse }\n" +
	"table.transactions td { border: 1px solid }\n" +
	"td.time { text-align: right }\n" +
	"td.id { text-align: right }\n" +
	"td.balance_amount, td.payment_amount { text-align: right }\n" +
	"td.balance_amount, th.balance { background-color: #e0ffff }\n" +
	"td.error { color: red }\n" +
	"td.success { color: green }\n" +
	"p.project { text-align: right }\n" +
	"\n" +
	"button#more { \n" +
	"    width: 100%;\n" +
	"    border: none;\n" +
	"    padding: 4px;\n" +
	"    background-color: #c0c0c0;\n" +
	"    color: black\n" +
	"}\n" +
	"\n" +
	".menu {\n" +
	"    background-color: black;\n" +
	"    color: white;\n" +
	"}\n" +
	"/* Style The Dropdown Button */\n" +
	".dropbtn {\n" +
	"    background-color: black;\n" +
	"    color: white;\n" +
	"    padding: 6pt;\n" +
	"    font-size: 12pt;\n" +
	"    border: 0pt;\n" +
	"    cursor: pointer;\n" +
	"    font-family: Liberation;\n" +
	"}\n" +
	"\n" +
	"/* The container <div> - needed to position the dropdown content" +
	" */\n" +
	".dropdown {\n" +
	"    position: relative;\n" +
	"    display: inline-block;\n" +
	"}\n" +
	"\n" +
	"/* Dropdown Content (Hidden by Default) */\n" +
	".dropdown-content {\n" +
	"    display: none;\n" +
	"    position: absolute;\n" +
	"    background-color: black;\n" +
	"    color: white;\n" +
	"    min-width: 160px;\n" +
	"    box-shadow: 0px 6pt 12pt 0px rgba(0,0,0,0.2);\n" +
	"    z-index: 1;\n" +
	"}\n" +
	"\n" +
	"/* Links inside the dropdown */\n" +
	".dropdown-content a {\n" +
	"    color: white;\n" +
	"    padding: 6pt;\n" +
	"    text-decoration: none;\n" +
	"    display: block;\n" +
	"}\n" +
	"\n" +
	"/* Change color of dropdown links on hover */\n" +
	".dropdown-content a:hover {\n" +
	"    background-color: #f0f0f0;\n" +
	"    color: black;\n" +
	"}\n" +
	"\n" +
	"/* Show the dropdown menu on hover */\n" +
	".dropdown:hover .dropdown-content {\n" +
	"    display: block;\n" +
	"}\n" +
	"\n" +
	"/* Change the background color of the dropdown button when the d" +
	"ropdown content is shown */\n" +
	".dropdown:hover .dropbtn {\n" +
	"    color: white;\n" +
	"    background-color: black;\n" +
	"}\n" +
	"")
const appjs = ("// Login page\n" +
	"\n" +
	"// login issues the login request.\n" +
	"function login() {\n" +
	"    $.ajax({\n" +
	"        method: \"POST\",\n" +
	"        data: JSON.stringify({\n" +
	"            User: $(\"#user\").val(),\n" +
	"            Password: $(\"#password\").val(),\n" +
	"        }),\n" +
	"        url: \"/v1/login\",\n" +
	"        contentType: \"application/json\",\n" +
	"        dataType: \"json\",\n" +
	"        success: function (status, data, jqxhr) {\n" +
	"            if (window.location.search == \"\") {\n" +
	"                window.location.href = \"/\";\n" +
	"            } else {\n" +
	"                window.location.href = window.location.search.su" +
	"bstring(1);\n" +
	"            }\n" +
	"        },\n" +
	"        error: function (jqxhr, error, exception) {\n" +
	"            $(\"#error\").text(jqxhr.responseText);\n" +
	"        },\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// Logout page\n" +
	"\n" +
	"// logout issues the logout request.\n" +
	"function logout() {\n" +
	"    $.ajax({\n" +
	"        method: \"POST\",\n" +
	"        url: \"/v1/logout\",\n" +
	"        contentType: \"application/json\",\n" +
	"        success: function () {\n" +
	"            window.location.href = \"/login.html\";\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"\n" +
	"// Transactions table\n" +
	"\n" +
	"// getTransactions extends the transactions table with older tra" +
	"nsactions\n" +
	"function transactionsMore() {\n" +
	"    $.ajax({\n" +
	"        url: \"/v1/transaction?\" + $.param({\n" +
	"            limit: 20,\n" +
	"            offset: $('table.transactions>tbody>tr').length\n" +
	"        }),\n" +
	"        dataType: \"json\",\n" +
	"        success: transactionsAddMore,\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false\n" +
	"}\n" +
	"\n" +
	"// addMoreTransactions adds a list of transactions to the end of" +
	" the transactions table\n" +
	"function transactionsAddMore(data) {\n" +
	"    var i;\n" +
	"    // Hide 'more' button if there's nothing to get\n" +
	"    if (data.length == 0) {\n" +
	"        $(\"#more\").hide()\n" +
	"        return\n" +
	"    }\n" +
	"    for (i = 0; i < data.length; i++) {\n" +
	"        $(\"table.transactions > tbody\").append(transactionToRo" +
	"w(data[i]));\n" +
	"    }\n" +
	"}\n" +
	"\n" +
	"// refreshTransactions extends the transactions table with older" +
	" transactions\n" +
	"function transactionsRefresh() {\n" +
	"    var td, id\n" +
	"    // Find the earliest transaction we know of\n" +
	"    td = $(\"table.transactions > tbody td.id\")\n" +
	"    if (td.length == 0) {\n" +
	"        transactionsMore()\n" +
	"        return\n" +
	"    }\n" +
	"    id = Number(td[0].innerText)\n" +
	"    $.ajax({\n" +
	"        url: \"/v1/transaction?\" + $.param({ after: id }),\n" +
	"        dataType: \"json\",\n" +
	"        success: transactionsAddNew,\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// addNewTransactions adds a list of transactions to the start o" +
	"f the transactions table\n" +
	"function transactionsAddNew(data) {\n" +
	"    var i;\n" +
	"    for (i = data.length - 1; i >= 0; i--) {\n" +
	"        $(\"table.transactions > tbody\").prepend(transactionToR" +
	"ow(data[i]));\n" +
	"    }\n" +
	"}\n" +
	"\n" +
	"// transactionToRow returns a <tr> element for a transaction.\n" +
	"function transactionToRow(transaction) {\n" +
	"    var tr, td, time, j\n" +
	"    tr = $(\"<tr>\");\n" +
	"    time = new Date(transaction[\"Time\"] * 1000)\n" +
	"    time = (time.getUTCFullYear() + \"-\" + pad(time.getUTCMonth" +
	"() + 1, 2) + \"-\" + pad(time.getUTCDate(), 2)\n" +
	"        + \" \" + pad(time.getUTCHours(), 2) + \":\" + pad(time." +
	"getUTCMinutes(), 2) + \":\" + pad(time.getUTCSeconds(), 2))\n" +
	"    td = $(\"<td>\").addClass(\"time\").text(time);\n" +
	"    tr.append(td)\n" +
	"    td = $(\"<td>\").addClass(\"id\").text(transaction[\"ID\"]);" +
	"\n" +
	"    tr.append(td)\n" +
	"    td = $(\"<td>\").addClass(\"user\").text(transaction[\"User\"" +
	"])\n" +
	"    tr.append(td)\n" +
	"    td = $(\"<td>\").addClass(\"description\").text(transaction[" +
	"\"Description\"])\n" +
	"    tr.append(td)\n" +
	"    td = $(\"<td>\").addClass(\"payment_amount\").text((transact" +
	"ion[\"Amount\"] / 100).toFixed(2))\n" +
	"    tr.append(td)\n" +
	"    td = $(\"<td>\").addClass(\"account\").text(transaction[\"Or" +
	"igin\"])\n" +
	"    tr.append(td)\n" +
	"    td = $(\"<td>\").addClass(\"account\").text(transaction[\"De" +
	"stination\"])\n" +
	"    tr.append(td)\n" +
	"    for (j = 0; j < accounts.length; j++) {\n" +
	"        td = $(\"<td>\").addClass(\"balance_amount\");\n" +
	"        if (accounts[j] == transaction[\"Origin\"]) {\n" +
	"            td.text((transaction[\"OriginBalanceAfter\"] / 100)." +
	"toFixed(2))\n" +
	"        } else if (accounts[j] == transaction[\"Destination\"]) " +
	"{\n" +
	"            td.text((transaction[\"DestinationBalanceAfter\"] / " +
	"100).toFixed(2));\n" +
	"        }\n" +
	"        tr.append(td);\n" +
	"    }\n" +
	"    return tr;\n" +
	"}\n" +
	"\n" +
	"// New transactions page\n" +
	"\n" +
	"// transactionsNew issues a transaction creation request.\n" +
	"function transactionsNew() {\n" +
	"    $.ajax({\n" +
	"        method: \"POST\",\n" +
	"        data: JSON.stringify({\n" +
	"            Token: $(\"#token\").val(),\n" +
	"            Origin: $(\"#origin\").val(),\n" +
	"            Destination: $(\"#destination\").val(),\n" +
	"            Description: $(\"#description\").val(),\n" +
	"            Amount: Math.round(Number($(\"#amount\").val()) * 10" +
	"0),\n" +
	"        }),\n" +
	"        url: \"/v1/transaction/\",\n" +
	"        contentType: \"application/json\",\n" +
	"        success: function () {\n" +
	"            transactionsRefresh();\n" +
	"            $(\"#success\").text(\"transaction created\");\n" +
	"            $(\"#description\").val(\"\")\n" +
	"            $(\"#amount\").val(\"\")\n" +
	"\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// Distribute page\n" +
	"\n" +
	"// distribute issues a transaction creation request.\n" +
	"function distribute() {\n" +
	"    $.ajax({\n" +
	"        method: \"POST\",\n" +
	"        data: JSON.stringify({\n" +
	"            Token: $(\"#token\").val(),\n" +
	"            Origin: $(\"#origin\").val(),\n" +
	"            Destinations: $(\"#destinations\").val(),\n" +
	"            Description: $(\"#description\").val(),\n" +
	"        }),\n" +
	"        url: \"/v1/distribute/\",\n" +
	"        contentType: \"application/json\",\n" +
	"        success: function () {\n" +
	"            transactionsRefresh();\n" +
	"            $(\"#success\").text(\"transaction created\");\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// New user page\n" +
	"\n" +
	"// newUser issues a user creation request.\n" +
	"function newUser() {\n" +
	"    var user;\n" +
	"    user = $(\"#user\").val();\n" +
	"    $.ajax({\n" +
	"        method: \"POST\",\n" +
	"        data: JSON.stringify({\n" +
	"            Token: $(\"#token\").val(),\n" +
	"            User: $(\"#user\").val(),\n" +
	"            Password: $(\"#password\").val(),\n" +
	"        }),\n" +
	"        url: \"/v1/user/\",\n" +
	"        contentType: \"application/json\",\n" +
	"        success: function () {\n" +
	"            $(\"#success\").text(\"user \" + user + \" created\"" +
	");\n" +
	"            // Clear the form. This seems a bit ugly but it does" +
	"n't make sense\n" +
	"            // to re-use usernames or passwords so I think it's " +
	"justifiable,\n" +
	"            // and it saves having to retrigger validation when " +
	"the user list\n" +
	"            // is updated!\n" +
	"            $(\"#user\").val(\"\");\n" +
	"            $(\"#password\").val(\"\");\n" +
	"            $(\"#password2\").val(\"\");\n" +
	"            updateUsers();\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// Delete user page\n" +
	"\n" +
	"// delUser issues a user deletion request.\n" +
	"function delUser() {\n" +
	"    var user;\n" +
	"    user = $(\"#user\").val();\n" +
	"    $.ajax({\n" +
	"        method: \"DELETE\",\n" +
	"        url: \"/v1/user/\" + user,\n" +
	"        success: function () {\n" +
	"            $(\"#success\").text(\"user \" + user + \" deleted\"" +
	");\n" +
	"            $(\"#user\").val(\"\");\n" +
	"            updateUsers();\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// New account page\n" +
	"\n" +
	"// newAccount issues a account creation request.\n" +
	"function newAccount() {\n" +
	"    var account;\n" +
	"    account = $(\"#account\").val();\n" +
	"    $.ajax({\n" +
	"        method: \"POST\",\n" +
	"        data: JSON.stringify({\n" +
	"            Token: $(\"#token\").val(),\n" +
	"            Account: account,\n" +
	"        }),\n" +
	"        url: \"/v1/account/\",\n" +
	"        contentType: \"application/json\",\n" +
	"        success: function () {\n" +
	"            $(\"#success\").text(\"account \" + account + \" cre" +
	"ated\");\n" +
	"            // Clear the form. See discussion in newUser.\n" +
	"            $(\"#account\").val(\"\");\n" +
	"            updateAccounts();\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// Delete account page\n" +
	"\n" +
	"// delAccount issues an account deletion request.\n" +
	"function delAccount() {\n" +
	"    var account;\n" +
	"    account = $(\"#account\").val();\n" +
	"    $.ajax({\n" +
	"        method: \"DELETE\",\n" +
	"        url: \"/v1/account/\" + account,\n" +
	"        success: function () {\n" +
	"            $(\"#success\").text(\"account \" + account + \" del" +
	"eted\");\n" +
	"            $(\"#account\").val(\"\");\n" +
	"            updateAccounts();\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// Change password page\n" +
	"\n" +
	"// changePassword issues a password change request\n" +
	"function changePassword() {\n" +
	"    var user;\n" +
	"    $.ajax({\n" +
	"        method: \"PUT\",\n" +
	"        data: JSON.stringify({\n" +
	"            Token: $(\"#token\").val(),\n" +
	"            Password: $(\"#password\").val(),\n" +
	"        }),\n" +
	"        url: \"/v1/user/\" + $(\"#user\").val() + \"/password\"," +
	"\n" +
	"        contentType: \"application/json\",\n" +
	"        success: function () {\n" +
	"            $(\"#success\").text(\"password changed\");\n" +
	"        },\n" +
	"        error: ajaxFailed,\n" +
	"    });\n" +
	"    return false;\n" +
	"}\n" +
	"\n" +
	"// Utilities\n" +
	"\n" +
	"// pad pads the integer n to d digits.\n" +
	"function pad(n, d) {\n" +
	"    return n.toString().padStart(d, \"0\");\n" +
	"}\n" +
	"\n" +
	"// amountRegexp matches valid currency amounts.\n" +
	"var amountRegexp = /^[0-9]+(\\.[0-9]{2})?$/;\n" +
	"\n" +
	"// validate validates form entries and adjusts the submit button" +
	".\n" +
	"function validate(container) {\n" +
	"    var valid;\n" +
	"    // Adjust the cooked transaction form.\n" +
	"    if ($(\"select#reason\").val() == \"house\") {\n" +
	"        $(\"#originRow\").hide();\n" +
	"        $(\"#origin\").val(config[\"houseAccount\"]); // TODO as" +
	"sumed to exist\n" +
	"        $(\"#origin\").removeClass(\"human\");\n" +
	"    } else if ($(\"select#reason\").val() == \"payback\") {\n" +
	"        $(\"#originRow\").show();\n" +
	"        $(\"#origin\").addClass(\"human\");\n" +
	"    }\n" +
	"    valid = true;\n" +
	"    container.find('td.error').text('');\n" +
	"    container.find('input,select').each(function (i, e) {\n" +
	"        var j, newpasswords, trouble;\n" +
	"        if (valid) {\n" +
	"            e = $(e)\n" +
	"            if (e.hasClass('nonempty') && e.val() == \"\") {\n" +
	"                trouble = \"must not be empty\";\n" +
	"                valid = false;\n" +
	"            }\n" +
	"            if (valid && e.hasClass(\"newuser\") && users.includ" +
	"es(e.val())) {\n" +
	"                trouble = \"user name in use\";\n" +
	"                valid = false;\n" +
	"            }\n" +
	"            if (valid && e.hasClass(\"newaccount\") && accounts." +
	"includes(e.val())) {\n" +
	"                trouble = \"account name in use\";\n" +
	"                valid = false;\n" +
	"            }\n" +
	"            if (valid && e.hasClass('currency') && !amountRegexp" +
	".exec(e.val())) {\n" +
	"                trouble = \"must be NN or NN.NN\";\n" +
	"                valid = false;\n" +
	"            }\n" +
	"            if (valid && e.hasClass('account') && e[0].id == \"d" +
	"estination\" && $(\"#origin\").val() == e.val()) {\n" +
	"                trouble = \"must not be the same as origin\";\n" +
	"                valid = false;\n" +
	"            }\n" +
	"            if (valid && e.hasClass('accounts')) {\n" +
	"                if (valid && e.val().length < 2) {\n" +
	"                    trouble = \"must distribute between multiple" +
	" accounts\";\n" +
	"                    valid = false;\n" +
	"                }\n" +
	"                if (valid && e.val().includes($(\"#origin\").val" +
	"())) {\n" +
	"                    // assumed to be destination\n" +
	"                    trouble = \"must not include origin (\" + $(" +
	"\"#origin\").val() + \")\";\n" +
	"                    valid = false;\n" +
	"                }\n" +
	"            }\n" +
	"            if (valid && e.hasClass('human') && e.val() == confi" +
	"g[\"houseAccount\"]) {\n" +
	"                trouble = \"must not be house account\";\n" +
	"                valid = false;\n" +
	"            }\n" +
	"            if (valid && e.hasClass('newpassword')) {\n" +
	"                newpasswords = []\n" +
	"                newpasswords = container.find(\"input.newpasswor" +
	"d\").each(function (i, p) {\n" +
	"                    if (e.val() != $(p).val()) {\n" +
	"                        trouble = \"passwords must match\";\n" +
	"                        valid = false;\n" +
	"                    }\n" +
	"                });\n" +
	"            }\n" +
	"            e.parent().next().text(trouble);\n" +
	"        }\n" +
	"    })\n" +
	"    container.find(\".submit\").prop(\"disabled\", !valid);\n" +
	"}\n" +
	"\n" +
	"// Error handling\n" +
	"\n" +
	"// ajaxFailed is called when a non-authentication error occurs\n" +
	"function ajaxFailed(jqxhr, error, exception) {\n" +
	"    if (jqxhr.status == 403) { // If logged out, bounce back to " +
	"login page\n" +
	"        if (window.location.pathname != \"/\") {\n" +
	"            window.location.href = \"/login.html?\" + window.loc" +
	"ation.pathname;\n" +
	"        } else {\n" +
	"            window.location.href = \"/login.html\";\n" +
	"        }\n" +
	"        return\n" +
	"    }\n" +
	"    $(\"#error\").text(jqxhr.responseText);\n" +
	"}\n" +
	"\n" +
	"// Initialization\n" +
	"// initialize gets configuration and sets up everything that nee" +
	"ds it.\n" +
	"function initialize() {\n" +
	"    var u, a, c;\n" +
	"    // Disable all form submission until ready\n" +
	"    $(\".submit\").prop(\"disabled\", true);\n" +
	"    // Get configuration &c concurrently\n" +
	"    u = $.ajax({\n" +
	"        url: \"/v1/user/\",\n" +
	"        dataType: \"json\",\n" +
	"    });\n" +
	"    a = $.ajax({\n" +
	"        url: \"/v1/account/\",\n" +
	"        dataType: \"json\",\n" +
	"    });\n" +
	"    c = $.ajax({\n" +
	"        url: \"/v1/config/\",\n" +
	"        dataType: \"json\",\n" +
	"    });\n" +
	"    $.when(u, a, c).done(function (ur, ar, cr) {\n" +
	"        var i, tr, select;\n" +
	"        config = cr[0];\n" +
	"        newUsers(ur[0])\n" +
	"        newAccounts(ar[0]);\n" +
	"        // Initialize the transactions table if present\n" +
	"        if ($(\"table.transactions\").length > 0) {\n" +
	"            // Populate the balance columns\n" +
	"            //\n" +
	"            // This only happens at page load and therefore does" +
	" not cope with\n" +
	"            // the set of accounts changing during the page's li" +
	"fetime.\n" +
	"            // I'm going to leave this as it is; the effort to f" +
	"ix it isn't\n" +
	"            // justified by the expected usage model.\n" +
	"            tr = $(\"table.transactions > thead > tr\");\n" +
	"            for (i = 0; i < accounts.length; i++) {\n" +
	"                tr.append($(\"<th>\").addClass(\"balance\").text" +
	"(accounts[i]));\n" +
	"            }\n" +
	"            // Populate the transactions table\n" +
	"            transactionsMore();\n" +
	"            // Refresh the transactions table regularly.\n" +
	"            // This has the effect of auto-bouncing you to /logi" +
	"n.html when things go wrong.\n" +
	"            // This isn't completely terrible (the login page wi" +
	"ll take you back to where\n" +
	"            // you were) but it's not very pretty. I'll leave it" +
	" for now.\n" +
	"            setInterval(transactionsRefresh, 10000);\n" +
	"            $(\"#more\").on(\"click\", transactionsMore);\n" +
	"        }\n" +
	"        initializeValidation();\n" +
	"        initializeForms();\n" +
	"    }).fail(ajaxFailed);\n" +
	"}\n" +
	"\n" +
	"function updateAccounts() {\n" +
	"    $.ajax({\n" +
	"        url: \"/v1/account/\",\n" +
	"        dataType: \"json\",\n" +
	"        success: newAccounts,\n" +
	"    });\n" +
	"}\n" +
	"\n" +
	"function newAccounts(a) {\n" +
	"    var i, select;\n" +
	"    accounts = a;\n" +
	"    // Populate the account dropdowns\n" +
	"    select = $(\"select.account,select.accounts\");\n" +
	"    select.find('option').remove()\n" +
	"    for (i = 0; i < accounts.length; i++) {\n" +
	"        select.append($(\"<option>\").text(accounts[i]));\n" +
	"    }\n" +
	"    // If there is a house account it should be the default tran" +
	"saction origin\n" +
	"    // for generic transactions; but we don't do this for the co" +
	"oked transaction\n" +
	"    // form.\n" +
	"    if (accounts.includes(config[\"houseAccount\"])) {\n" +
	"        if (!$(\"select#origin\").hasClass(\"human\")) {\n" +
	"            $(\"select#origin\").val(config[\"houseAccount\"]);\n" +
	"        }\n" +
	"    }\n" +
	"}\n" +
	"\n" +
	"function updateUsers() {\n" +
	"    $.ajax({\n" +
	"        url: \"/v1/user/\",\n" +
	"        dataType: \"json\",\n" +
	"        success: newUsers,\n" +
	"    });\n" +
	"}\n" +
	"\n" +
	"function newUsers(u) {\n" +
	"    var i, select;\n" +
	"    users = u;\n" +
	"    // Populate the user dropdowns\n" +
	"    select = $(\"select.user,select.users\");\n" +
	"    select.find('option').remove()\n" +
	"    for (i = 0; i < users.length; i++) {\n" +
	"        select.append($(\"<option>\").text(users[i]));\n" +
	"    }\n" +
	"}\n" +
	"\n" +
	"// initializeValidation sets up form validation logic\n" +
	"// and does the initial validation of the (mostly empty)\n" +
	"// forms.\n" +
	"function initializeValidation() {\n" +
	"    // Whenever any form is modified...\n" +
	"    $(\"form\").each(function (i, f) {\n" +
	"        $(f).find(\"input,select\").on(\"input change\", functio" +
	"n () {\n" +
	"            // ...clear the error & success indicators\n" +
	"            $(\"#error\").text(\"\");\n" +
	"            $(\"#success\").text(\"\");\n" +
	"            // ...and revalidate the form\n" +
	"            validate($(f));\n" +
	"        });\n" +
	"    });\n" +
	"    // Attach requests to form submission\n" +
	"    $(\"form#login\").on(\"submit\", login);\n" +
	"    $(\"form#newTransaction\").on(\"submit\", transactionsNew);\n" +
	"    $(\"form#cookedTransaction\").on(\"submit\", transactionsNew" +
	");\n" +
	"    $(\"form#distribute\").on(\"submit\", distribute);\n" +
	"    $(\"form#newuser\").on(\"submit\", newUser);\n" +
	"    $(\"form#deluser\").on(\"submit\", delUser);\n" +
	"    $(\"form#newaccount\").on(\"submit\", newAccount);\n" +
	"    $(\"form#delaccount\").on(\"submit\", delAccount);\n" +
	"    $(\"form#changepass\").on(\"submit\", changePassword);\n" +
	"    // Initial validation of forms\n" +
	"    $(\"form\").each(function (i, f) {\n" +
	"        validate($(f));\n" +
	"    });\n" +
	"}\n" +
	"\n" +
	"// initializeForms sets initial values for some of the forms.\n" +
	"function initializeForms() {\n" +
	"    // Default human is logged in user, if they have an account\n" +
	"    $(\"select.human\").each(function (i, f) {\n" +
	"        user = $(\"input#user\").val()\n" +
	"        if (accounts.includes(user)) {\n" +
	"            $(f).val(user)\n" +
	"        } else {\n" +
	"            $(f).val(\"\")\n" +
	"        }\n" +
	"    });\n" +
	"}\n" +
	"\n" +
	"$(document).ready(function () {\n" +
	"    // Login and logout are special.\n" +
	"    if (window.location.pathname == \"/logout.html\") {\n" +
	"        logout();\n" +
	"    } else if (window.location.pathname == \"/login.html\") {\n" +
	"        initializeValidation();\n" +
	"    } else {\n" +
	"        initialize();\n" +
	"    }\n" +
	"    // Logout is ready immediately.\n" +
	"    $(\"a#logout\").on(\"click\", logout);\n" +
	"});\n" +
	"")
const delaccounthtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Delete Account ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Delete Account</h1>\n" +
	"\n" +
	"    <form id=delaccount>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>User</td>\n" +
	"                <td>\n" +
	"                    <select name=\"account\" class=\"account\" i" +
	"d=account></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input type=submit value=Delete class=submit" +
	">\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"                <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>\n" +
	"")
const deluserhtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Delete User ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Delete User</h1>\n" +
	"\n" +
	"    <form id=deluser>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>User</td>\n" +
	"                <td>\n" +
	"                    <select name=\"user\" class=\"user\" id=user" +
	"></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input type=submit value=Delete class=submit" +
	">\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"                <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>\n" +
	"")
const distributehtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Distribute ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\"\n" +
	"        integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8Qtmk" +
	"MRdAu8=\" crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Distribute</h1>\n" +
	"\n" +
	"    <form id=distribute>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>Description</td>\n" +
	"                <td>\n" +
	"                    <input name=\"description\" type=text value=" +
	"\"distribution\" id=description class=nonempty>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Origin</td>\n" +
	"                <td>\n" +
	"                    <select name=\"origin\" class=account id=ori" +
	"gin></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Destination</td>\n" +
	"                <td>\n" +
	"                    <select multiple name=\"destinations\" class" +
	"=accounts id=destinations>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input name=\"submit\" type=submit value=\"D" +
	"istribute\" class=submit>\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td colspan=3 id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"\n" +
	"    <div class=transactions>\n" +
	"        <table class=transactions>\n" +
	"            <thead>\n" +
	"                <tr>\n" +
	"                    <th>Time</th>\n" +
	"                    <th>ID</th>\n" +
	"                    <th>User</th>\n" +
	"                    <th>Description</th>\n" +
	"                    <th>Amount</th>\n" +
	"                    <th>From</th>\n" +
	"                    <th>To</th>\n" +
	"                </tr>\n" +
	"            </thead>\n" +
	"            <tbody>\n" +
	"            </tbody>\n" +
	"        </table>\n" +
	"        <button id=more>More...</button>\n" +
	"    </div>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>\n" +
	"")
const indexhtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Transactions ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\"\n" +
	"        integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8Qtmk" +
	"MRdAu8=\" crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Transactions</h1>\n" +
	"\n" +
	"    <form id=cookedTransaction>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>What happened</td>\n" +
	"                <td>\n" +
	"                    <select name=\"reason\" id=reason>\n" +
	"                        <option value=house selected>Bought a sh" +
	"ared resource</option>\n" +
	"                        <option value=payback>Repaid another par" +
	"ty</option>\n" +
	"                    </select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Who paid</td>\n" +
	"                <td>\n" +
	"                    <select name=\"destination\" class=\"account" +
	" human\" id=destination></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr id=originRow>\n" +
	"                <td>Who was paid</td>\n" +
	"                <td>\n" +
	"                    <select name=\"origin\" class=\"account\" id" +
	"=origin></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Description</td>\n" +
	"                <td>\n" +
	"                    <input name=\"description\" type=text value=" +
	"\"\" id=description class=nonempty>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Amount</td>\n" +
	"                <td>\n" +
	"                    <input name=\"amount\" type=text value=\"\" " +
	"id=amount class=currency>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input name=\"submit\" type=submit value=\"G" +
	"o\" class=submit>\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"                <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"        <input name=user type=hidden value=\"{{.User}}\" id=user" +
	">\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"\n" +
	"    <div class=transactions>\n" +
	"        <table class=transactions>\n" +
	"            <thead>\n" +
	"                <tr>\n" +
	"                    <th>Time</th>\n" +
	"                    <th>ID</th>\n" +
	"                    <th>User</th>\n" +
	"                    <th>Description</th>\n" +
	"                    <th>Amount</th>\n" +
	"                    <th>From</th>\n" +
	"                    <th>To</th>\n" +
	"                </tr>\n" +
	"            </thead>\n" +
	"            <tbody>\n" +
	"            </tbody>\n" +
	"        </table>\n" +
	"        <button id=more>More...</button>\n" +
	"    </div>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>\n" +
	"")
const loginhtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Login ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Login</h1>\n" +
	"\n" +
	"    <form id=login>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>User name</td>\n" +
	"                <td>\n" +
	"                    <input id=user type=text autocomplete=userna" +
	"me value=\"\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Password</td>\n" +
	"                <td>\n" +
	"                    <input id=password type=password autocomplet" +
	"e=\"current-password\" value=\"\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input type=submit value=Login class=submit>" +
	"\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>\n" +
	"")
const logouthtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Logout ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Logout</h1>\n" +
	"\n" +
	"    <p id=error class=error></p>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>")
const newaccounthtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>New Account ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: New Account</h1>\n" +
	"\n" +
	"    <form id=newaccount>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>Account name</td>\n" +
	"                <td>\n" +
	"                    <input id=account type=text value=\"\" class" +
	"=\"nonempty newaccount\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr></tr>\n" +
	"            <td></td>\n" +
	"            <td>\n" +
	"                <input type=submit value=Create class=submit>\n" +
	"            </td>\n" +
	"            <td id=error class=error></td>\n" +
	"            <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>")
const newuserhtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>New User ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: New User</h1>\n" +
	"\n" +
	"    <form id=newuser>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>User name</td>\n" +
	"                <td>\n" +
	"                    <input id=user type=text value=\"\" class=\"" +
	"nonempty newuser\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Password</td>\n" +
	"                <td>\n" +
	"                    <input id=password type=password autocomplet" +
	"e=\"new-password\" value=\"\" class=\"nonempty newpassword\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Confirm password</td>\n" +
	"                <td>\n" +
	"                    <input id=password2 type=password autocomple" +
	"te=\"new-password\" value=\"\" class=\"nonempty newpassword\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input type=submit value=Create class=submit" +
	">\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"                <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>")
const passwordhtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>Change Password ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\" " +
	"integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=\"" +
	"\n" +
	"        crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: Change Password</h1>\n" +
	"\n" +
	"    <form id=changepass>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>New password</td>\n" +
	"                <td>\n" +
	"                    <input id=password type=password autocomplet" +
	"e=\"new-password\" value=\"\" class=\"nonempty newpassword\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Confirm password</td>\n" +
	"                <td>\n" +
	"                    <input id=password2 type=password autocomple" +
	"te=\"new-password\" value=\"\" class=\"nonempty newpassword\">\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input type=submit value=Change class=submit" +
	">\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"                <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"        <input name=user type=hidden autocomplete=\"username\" v" +
	"alue=\"{{.User}}\" id=user>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>")
const transactionhtml = ("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">\n" +
	"<html>\n" +
	"\n" +
	"<head>\n" +
	"    <title>New Transaction ({{.Title}})</title>\n" +
	"    <script src=\"https://code.jquery.com/jquery-3.3.1.min.js\"\n" +
	"        integrity=\"sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8Qtmk" +
	"MRdAu8=\" crossorigin=\"anonymous\"></script>\n" +
	"    <script type=\"text/javascript\" src=\"/app.js\"></script>\n" +
	"    <link rel=StyleSheet type=\"text/css\" href=\"/app.css\">\n" +
	"</head>\n" +
	"\n" +
	"<body>\n" +
	"    <div class=menu>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Transactions</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/\">Transactions</a>\n" +
	"                <a href=\"/transaction.html\">New Transaction</a" +
	">\n" +
	"                <a href=\"/distribute.html\">Distribute</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Accounts</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newaccount.html\">New Account</a>\n" +
	"                <a href=\"/delaccount.html\">Delete Account</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"        <div class=dropdown>\n" +
	"            <button class=dropbtn>Users</button>\n" +
	"            <div class=\"dropdown-content\">\n" +
	"                <a href=\"/newuser.html\">New User</a>\n" +
	"                <a href=\"/deluser.html\">Delete User</a>\n" +
	"                <a href=\"/password.html\">Change Password</a>\n" +
	"                <a id=logout href=\"/logout.html\">Logout</a>\n" +
	"            </div>\n" +
	"        </div>\n" +
	"    </div>\n" +
	"\n" +
	"    <h1>{{.Title}}: New Transaction</h1>\n" +
	"\n" +
	"    <form id=newTransaction>\n" +
	"        <table class=form>\n" +
	"            <tr>\n" +
	"                <td>Description</td>\n" +
	"                <td>\n" +
	"                    <input name=\"description\" type=text value=" +
	"\"\" id=description class=nonempty>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Origin</td>\n" +
	"                <td>\n" +
	"                    <select name=\"origin\" class=\"account\" id" +
	"=origin></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Destination</td>\n" +
	"                <td>\n" +
	"                    <select name=\"destination\" class=account i" +
	"d=destination></select>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td>Amount</td>\n" +
	"                <td>\n" +
	"                    <input name=\"amount\" type=text value=\"\" " +
	"id=amount class=currency>\n" +
	"                </td>\n" +
	"                <td class=error></td>\n" +
	"            </tr>\n" +
	"            <tr>\n" +
	"                <td></td>\n" +
	"                <td>\n" +
	"                    <input name=\"submit\" type=submit value=\"N" +
	"ew Transaction\" class=submit>\n" +
	"                </td>\n" +
	"                <td id=error class=error></td>\n" +
	"                <td id=success class=success></td>\n" +
	"            </tr>\n" +
	"        </table>\n" +
	"        <input name=token type=hidden value=\"{{.Token}}\" id=to" +
	"ken>\n" +
	"    </form>\n" +
	"\n" +
	"    <hr>\n" +
	"\n" +
	"    <div class=transactions>\n" +
	"        <table class=transactions>\n" +
	"            <thead>\n" +
	"                <tr>\n" +
	"                    <th>Time</th>\n" +
	"                    <th>ID</th>\n" +
	"                    <th>User</th>\n" +
	"                    <th>Description</th>\n" +
	"                    <th>Amount</th>\n" +
	"                    <th>From</th>\n" +
	"                    <th>To</th>\n" +
	"                </tr>\n" +
	"            </thead>\n" +
	"            <tbody>\n" +
	"            </tbody>\n" +
	"        </table>\n" +
	"        <button id=more>More...</button>\n" +
	"    </div>\n" +
	"\n" +
	"    <hr>\n" +
	"    <p class=project><a href=\"https://github.com/ewxrjk/bank\">" +
	"Bank</a> ({{.Version}})</p>\n" +
	"</body>\n" +
	"\n" +
	"</html>\n" +
	"")

var embedContent = map[string]string{
	"app.css":          appcss,
	"app.js":           appjs,
	"delaccount.html":  delaccounthtml,
	"deluser.html":     deluserhtml,
	"distribute.html":  distributehtml,
	"index.html":       indexhtml,
	"login.html":       loginhtml,
	"logout.html":      logouthtml,
	"newaccount.html":  newaccounthtml,
	"newuser.html":     newuserhtml,
	"password.html":    passwordhtml,
	"transaction.html": transactionhtml,
}
var embedType = map[string]string{
	"app.css":          "text/css",
	"app.js":           "application/javascript",
	"delaccount.html":  "text/html",
	"deluser.html":     "text/html",
	"distribute.html":  "text/html",
	"index.html":       "text/html",
	"login.html":       "text/html",
	"logout.html":      "text/html",
	"newaccount.html":  "text/html",
	"newuser.html":     "text/html",
	"password.html":    "text/html",
	"transaction.html": "text/html",
}

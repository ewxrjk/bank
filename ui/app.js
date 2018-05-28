// Login page

// login issues the login request.
function login() {
    // TODO https://www.chromium.org/developers/design-documents/create-amazing-password-forms for advice
    $.ajax({
        method: "POST",
        data: JSON.stringify({
            User: $("#user").val(),
            Password: $("#password").val(),
        }),
        url: "/v1/login",
        contentType: "application/json",
        dataType: "json",
        success: function (status, data, jqxhr) {
            if (window.location.search == "") {
                window.location.href = "/";
            } else {
                window.location.href = window.location.search.substring(1);
            }
        },
        error: function (jqxhr, error, exception) {
            $("#error").text(jqxhr.responseText);
        },
    });
    return false;
}

// Logout page

// logout issues the logout request.
function logout() {
    $.ajax({
        method: "POST",
        url: "/v1/logout",
        contentType: "application/json",
        success: function () {
            window.location.href = "/login.html";
        },
        error: ajaxFailed,
    });
    return false;
}


// Transactions table

// getTransactions extends the transactions table with older transactions
function transactionsMore() {
    $.ajax({
        url: "/v1/transaction?" + $.param({
            limit: 20,
            offset: $('table.transactions>tbody>tr').length
        }),
        dataType: "json",
        success: transactionsAddMore,
        error: ajaxFailed,
    });
    return false
}

// addMoreTransactions adds a list of transactions to the end of the transactions table
function transactionsAddMore(data) {
    var i;
    // Hide 'more' button if there's nothing to get
    if (data.length == 0) {
        $("#more").hide()
        return
    }
    for (i = 0; i < data.length; i++) {
        $("table.transactions > tbody").append(transactionToRow(data[i]));
    }
}

// refreshTransactions extends the transactions table with older transactions
function transactionsRefresh() {
    var td, id
    // Find the earliest transaction we know of
    td = $("table.transactions > tbody td.id")
    if (td.length == 0) {
        transactionsMore()
        return
    }
    id = Number(td[0].innerText)
    $.ajax({
        url: "/v1/transaction?" + $.param({ after: id }),
        dataType: "json",
        success: transactionsAddNew,
        error: ajaxFailed,
    });
    return false;
}

// addNewTransactions adds a list of transactions to the start of the transactions table
function transactionsAddNew(data) {
    var i;
    for (i = data.length - 1; i >= 0; i--) {
        $("table.transactions > tbody").prepend(transactionToRow(data[i]));
    }
}

// transactionToRow returns a <tr> element for a transaction.
function transactionToRow(transaction) {
    var tr, td, time, j
    tr = $("<tr>");
    time = new Date(transaction["Time"] * 1000)
    time = (time.getUTCFullYear() + "-" + pad(time.getUTCMonth(), 2) + "-" + pad(time.getUTCDate(), 2)
        + " " + pad(time.getUTCHours(), 2) + ":" + pad(time.getUTCMinutes(), 2) + ":" + pad(time.getUTCSeconds(), 2))
    td = $("<td>").addClass("time").text(time);
    tr.append(td)
    td = $("<td>").addClass("id").text(transaction["ID"]);
    tr.append(td)
    td = $("<td>").addClass("user").text(transaction["User"])
    tr.append(td)
    td = $("<td>").addClass("description").text(transaction["Description"])
    tr.append(td)
    td = $("<td>").addClass("payment_amount").text((transaction["Amount"] / 100).toFixed(2))
    tr.append(td)
    td = $("<td>").addClass("account").text(transaction["Origin"])
    tr.append(td)
    td = $("<td>").addClass("account").text(transaction["Destination"])
    tr.append(td)
    for (j = 0; j < accounts.length; j++) {
        td = $("<td>").addClass("balance_amount");
        if (accounts[j] == transaction["Origin"]) {
            td.text((transaction["OriginBalanceAfter"] / 100).toFixed(2))
        } else if (accounts[j] == transaction["Destination"]) {
            td.text((transaction["DestinationBalanceAfter"] / 100).toFixed(2));
        }
        tr.append(td);
    }
    return tr;
}

// transactionsInit sets up the transactions table on page load
// TODO doesn't cope with account list changing!
function transactionsInit() {
    // Populate the balance columns and account dropdowns
    var i, tr, select;
    tr = $("table.transactions > thead > tr");
    select = $("select.account,select.accounts");
    for (i = 0; i < accounts.length; i++) {
        tr.append($("<th>").addClass("balance").text(accounts[i]));
        select.append($("<option>").text(accounts[i]));
    }
    // Populate the transactions table
    transactionsMore();
    // Refresh the transactions table regularly
    setInterval(transactionsRefresh, 10000);
    // TODO this has the effect of auto-bouncing you to /login.html when things go wrong
    $("#more").on("click", transactionsMore);
}

// New transactions page

// newTransaction issues a transaction creation request.
function transactionsNew() {
    $.ajax({
        method: "POST",
        data: JSON.stringify({
            Token: $("#token").val(),
            Origin: $("#origin").val(),
            Destination: $("#destination").val(),
            Description: $("#description").val(),
            Amount: Number($("#amount").val()) * 100,
        }),
        url: "/v1/transaction/",
        contentType: "application/json",
        success: function () {
            transactionsRefresh();
            $("#success").text("transaction created");
        },
        error: ajaxFailed,
    });
    return false;
}

// Distribute page

// distribute issues a transaction creation request.
function distribute() {
    $.ajax({
        method: "POST",
        data: JSON.stringify({
            Token: $("#token").val(),
            Origin: $("#origin").val(),
            Destinations: $("#destinations").val(),
            Description: $("#description").val(),
        }),
        url: "/v1/distribute/",
        contentType: "application/json",
        success: function () {
            transactionsRefresh();
            $("#success").text("transaction created");
        },
        error: ajaxFailed,
    });
    return false;
}

// New user page

// newUser issues a user creation request.
function newUser() {
    var user;
    user = $("#user").val();
    $.ajax({
        method: "POST",
        data: JSON.stringify({
            Token: $("#token").val(),
            User: $("#user").val(),
            Password: $("#password").val(),
        }),
        url: "/v1/user/",
        contentType: "application/json",
        success: function () {
            $("#success").text("user " + user + " created");
            updateUsers();
            // Clear the form. This seems a bit ugly but it doesn't make sense
            // to re-use usernames or passwords so I think it's justifiable,
            // and it saves having to retrigger validation when the user list
            // is updated!
            $("#user").val("");
            $("#password").val("");
            $("#password2").val("");
        },
        error: ajaxFailed,
    });
    return false;
}

// New account page

// newAccount issues a account creation request.
function newAccount() {
    var account;
    account = $("#account").val();
    $.ajax({
        method: "POST",
        data: JSON.stringify({
            Token: $("#token").val(),
            Account: account,
        }),
        url: "/v1/account/",
        contentType: "application/json",
        success: function () {
            $("#success").text("account " + account + " created");
            updateAccounts();
            // Clear the form. See discussion in newUser.
            $("#account").val("");
        },
        error: ajaxFailed,
    });
    return false;
}

// Change password page

// changePassword issues a password change request
function changePassword() {
    var user;
    $.ajax({
        method: "POST",
        data: JSON.stringify({
            Token: $("#token").val(),
            Password: $("#password").val(),
        }),
        url: "/v1/user/password",
        contentType: "application/json",
        success: function () {
            $("#success").text("password changed");
        },
        error: ajaxFailed,
    });
    return false;
}

// Utilities

// pad pads the integer n to d digits.
function pad(n, d) {
    return n.toString().padStart(d, "0");
}

// amountRegexp matches valid currency amounts.
var amountRegexp = /^[0-9]+(\.[0-9]{2})?$/;

// validate validates form entries and adjusts the submit button.
function validate(container) {
    var valid;
    valid = true;
    container.find('td.error').text('');
    container.find('input,select').each(function (i, e) {
        var j, newpasswords, trouble;
        if (valid) {
            e = $(e)
            if (e.hasClass('nonempty') && e.val() == "") {
                trouble = "must not be empty";
                valid = false;
            }
            if (valid && e.hasClass("newuser") && users.includes(e.val())) {
                trouble = "user name in use";
                valid = false;
            }
            if (valid && e.hasClass("newaccount") && accounts.includes(e.val())) {
                trouble = "account name in use";
                valid = false;
            }
            if (valid && e.hasClass('currency') && !amountRegexp.exec(e.val())) {
                trouble = "must be NN or NN.NN";
                valid = false;
            }
            if (valid && e.hasClass('account') && e[0].id == "destination" && $("#origin").val() == e.val()) {
                trouble = "must not be the same as origin";
                valid = false;
            }
            if (valid && e.hasClass('accounts')) {
                if (valid && e.val().length < 2) {
                    trouble = "must distribute between multiple accounts";
                    valid = false;
                }
                if (valid && e.val().includes($("#origin").val())) {
                    // assumed to be destination
                    trouble = "must not include origin (" + $("#origin").val() + ")";
                    valid = false;
                }
            }
            if (valid && e.hasClass('newpassword')) {
                newpasswords = []
                newpasswords = container.find("input.newpassword").each(function (i, p) {
                    if (e.val() != $(p).val()) {
                        trouble = "passwords must match";
                        valid = false;
                    }
                });
            }
            e.parent().next().text(trouble);
        }
    })
    container.find(".submit").prop("disabled", !valid);
}

// Error handling

// ajaxFailed is called when a non-authentication error occurs
function ajaxFailed(jqxhr, error, exception) {
    if (jqxhr.status == 403) { // If logged out, bounce back to login page
        if (window.location.pathname != "/") {
            window.location.href = "/login.html?" + window.location.pathname;
        } else {
            window.location.href = "/login.html";
        }
        return
    }
    $("#error").text(jqxhr.responseText);
}

// Initialization

// initialize gets configuration and sets up everything that needs it.
function initialize() {
    var u, a, c;
    // Disable all form submission until ready
    $(".submit").prop("disabled", true);
    // Get configuration &c concurrently
    u = $.ajax({
        url: "/v1/user/",
        dataType: "json",
    });
    a = $.ajax({
        url: "/v1/account/",
        dataType: "json",
    });
    c = $.ajax({
        url: "/v1/config/",
        dataType: "json",
    });
    $.when(u, a, c).done(function (ur, ar, cr) {
        users = ur[0];
        accounts = ar[0]
        config = cr[0];
        // Initialize the transactions table if present
        if ($("table.transactions").length > 0) {
            transactionsInit()
        }
        // If there is a house account it should be the default transaction origin
        if (accounts.includes(config["houseAccount"])) {
            $("select#origin").val(config["houseAccount"]);
        }
        initializeValidation();
    }).fail(ajaxFailed);
}

// initializeValidation sets up form validation logic
// and does the initial validation of the (mostly empty)
// forms.
function initializeValidation() {
    // Whenever any form is modified...
    $("form").each(function (i, f) {
        $(f).find("input,select").on("input", function () {
            // ...clear the error & success indicators
            $("#error").text("");
            $("#success").text("");
            // ...and revalidate the form
            validate($(f));
        });
    });
    // Attach requests to form submission
    $("form#login").on("submit", login);
    $("form#newTransaction").on("submit", transactionsNew);
    $("form#distribute").on("submit", distribute);
    $("form#newuser").on("submit", newUser);
    $("form#newaccount").on("submit", newAccount);
    $("form#changepass").on("submit", changePassword);
    // Initial validation of forms
    $("form").each(function (i, f) {
        validate($(f));
    });
}

$(document).ready(function () {
    // Login and logout are special.
    if (window.location.pathname == "/logout.html") {
        logout();
    } else if (window.location.pathname == "/login.html") {
        initializeValidation();
    } else {
        initialize();
    }
    // Logout is ready immediately.
    $("a#logout").on("click", logout);
});

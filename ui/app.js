// Login page

// login issues the login request.
function login() {
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
            Amount: Math.round(Number($("#amount").val()) * 100),
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
            // Clear the form. This seems a bit ugly but it doesn't make sense
            // to re-use usernames or passwords so I think it's justifiable,
            // and it saves having to retrigger validation when the user list
            // is updated!
            $("#user").val("");
            $("#password").val("");
            $("#password2").val("");
            updateUsers();
        },
        error: ajaxFailed,
    });
    return false;
}

// Delete user page

// delUser issues a user deletion request.
function delUser() {
    var user;
    user = $("#user").val();
    $.ajax({
        method: "DELETE",
        url: "/v1/user/" + user,
        success: function () {
            $("#success").text("user " + user + " deleted");
            $("#user").val("");
            updateUsers();
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
            // Clear the form. See discussion in newUser.
            $("#account").val("");
            updateAccounts();
        },
        error: ajaxFailed,
    });
    return false;
}

// Delete account page

// delAccount issues an account deletion request.
function delAccount() {
    var account;
    account = $("#account").val();
    $.ajax({
        method: "DELETE",
        url: "/v1/account/" + account,
        success: function () {
            $("#success").text("account " + account + " deleted");
            $("#account").val("");
            updateAccounts();
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
        method: "PUT",
        data: JSON.stringify({
            Token: $("#token").val(),
            Password: $("#password").val(),
        }),
        url: "/v1/user/" + $("#user").val() + "/password",
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
    // Adjust the cooked transaction form.
    if ($("select#reason").val() == "house") {
        $("#originRow").hide();
        $("#origin").val(config["houseAccount"]); // TODO assumed to exist
        $("#origin").removeClass("human");
    } else if ($("select#reason").val() == "payback") {
        $("#originRow").show();
        $("#origin").addClass("human");
    }
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
            if (valid && e.hasClass('human') && e.val() == config["houseAccount"]) {
                trouble = "must not be house account";
                valid = false;
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
        var i, tr, select;
        config = cr[0];
        newUsers(ur[0])
        newAccounts(ar[0]);
        // Initialize the transactions table if present
        if ($("table.transactions").length > 0) {
            // Populate the balance columns
            //
            // This only happens at page load and therefore does not cope with
            // the set of accounts changing during the page's lifetime.
            // I'm going to leave this as it is; the effort to fix it isn't
            // justified by the expected usage model.
            tr = $("table.transactions > thead > tr");
            for (i = 0; i < accounts.length; i++) {
                tr.append($("<th>").addClass("balance").text(accounts[i]));
            }
            // Populate the transactions table
            transactionsMore();
            // Refresh the transactions table regularly.
            // This has the effect of auto-bouncing you to /login.html when things go wrong.
            // This isn't completely terrible (the login page will take you back to where
            // you were) but it's not very pretty. I'll leave it for now.
            setInterval(transactionsRefresh, 10000);
            $("#more").on("click", transactionsMore);
        }
        initializeValidation();
        initializeForms();
    }).fail(ajaxFailed);
}

function updateAccounts() {
    $.ajax({
        url: "/v1/account/",
        dataType: "json",
        success: newAccounts,
    });
}

function newAccounts(a) {
    var i, select;
    accounts = a;
    // Populate the account dropdowns
    select = $("select.account,select.accounts");
    select.find('option').remove()
    for (i = 0; i < accounts.length; i++) {
        select.append($("<option>").text(accounts[i]));
    }
    // If there is a house account it should be the default transaction origin
    // for generic transactions; but we don't do this for the cooked transaction
    // form.
    if (accounts.includes(config["houseAccount"])) {
        if (!$("select#origin").hasClass("human")) {
            $("select#origin").val(config["houseAccount"]);
        }
    }
}

function updateUsers() {
    $.ajax({
        url: "/v1/user/",
        dataType: "json",
        success: newUsers,
    });
}

function newUsers(u) {
    var i, select;
    users = u;
    // Populate the user dropdowns
    select = $("select.user,select.users");
    select.find('option').remove()
    for (i = 0; i < users.length; i++) {
        select.append($("<option>").text(users[i]));
    }
}

// initializeValidation sets up form validation logic
// and does the initial validation of the (mostly empty)
// forms.
function initializeValidation() {
    // Whenever any form is modified...
    $("form").each(function (i, f) {
        $(f).find("input,select").on("input change", function () {
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
    $("form#cookedTransaction").on("submit", transactionsNew);
    $("form#distribute").on("submit", distribute);
    $("form#newuser").on("submit", newUser);
    $("form#deluser").on("submit", delUser);
    $("form#newaccount").on("submit", newAccount);
    $("form#delaccount").on("submit", delAccount);
    $("form#changepass").on("submit", changePassword);
    // Initial validation of forms
    $("form").each(function (i, f) {
        validate($(f));
    });
}

// initializeForms sets initial values for some of the forms.
function initializeForms() {
    // Default human is logged in user, if they have an account
    $("select.human").each(function (i, f) {
        user = $("input#user").val()
        if (accounts.includes(user)) {
            $(f).val(user)
        } else {
            $(f).val("")
        }
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

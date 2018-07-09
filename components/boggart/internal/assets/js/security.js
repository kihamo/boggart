$(document).ready(function () {
    window.securityChangeStatus = function(status) {
        $.post('/boggart/security/', {'status': status}, function() {
            window.location.reload();
        });
    };
});


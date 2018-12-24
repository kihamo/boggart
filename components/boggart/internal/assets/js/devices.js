$(document).ready(function () {
    var tableDevices = $('#devices table')
        .DataTable({
            pageLength: 50,
            language: {
                url: '/dashboard/datatables/i18n.json'
            },
            ajax: {
                url: '/boggart/devices/?entity=devices',
                dataSrc: 'data'
            },
            columns: [
                {
                    data: 'types',
                    render: function (types) {
                        var content = '';

                        for (var i in types) {
                            content += '<span class="label label-success">' + types[i] + '</span> ';
                        }

                        return content;
                    }
                },
                {
                    data: 'status',
                    render: function (status) {
                        switch(status.toLowerCase()) {
                            case 'online':
                                return '<span class="label label-success">' + status + '</span>';

                            case 'offline':
                                return '<span class="label label-danger">' + status + '</span>';

                            case 'removing':
                                return '<span class="label label-info">' + status + '</span>';

                            case 'removed':
                                return '<span class="label label-warning">' + status + '</span>';

                            default:
                                return '<span class="label label-default">' + status + '</span>';
                        }
                    }
                },
                {
                    data: 'tasks',
                    render: function (tasks) {
                        return tasks.length;
                    }
                },
                {
                    data: 'mqtt_topics',
                    render: function (topics) {
                        return topics.length;
                    }
                },
                {
                    data: 'mqtt_subscribers',
                    render: function (subscribers) {
                        return subscribers.length;
                    }
                },
                {
                    data: 'serial_number'
                },
                {
                    data: 'id'
                },
                {
                    data: 'description'
                },
                {
                    data: 'tasks',
                    render: function (tasks) {
                        var content = '';

                        for (var i in tasks) {
                            content += '<span class="label label-warning">' + tasks[i] + '</span> ';
                        }

                        return content;
                    }
                },
                {
                    data: 'mqtt_topics',
                    render: function (topics) {
                        var content = '';

                        for (var i in topics) {
                            content += '<span class="label label-primary">' + topics[i] + '</span> ';
                        }

                        return content;
                    }
                },
                {
                    data: 'mqtt_subscribers',
                    render: function (subscribers) {
                        var content = '';

                        for (var i in subscribers) {
                            content += '<span class="label label-info">' + subscribers[i] + '</span> ';
                        }

                        return content;
                    }
                }
            ]
        });

    var tableListeners = $('#listeners table')
        .DataTable({
            language: {
                url: '/dashboard/datatables/i18n.json'
            },
            ajax: {
                url: '/boggart/devices/?entity=listeners',
                dataSrc: 'data'
            },
            columns: [
                {
                    data: 'name'
                },
                {
                    data: 'id'
                },
                {
                    data: 'fires'
                },
                {
                    data: 'fire_first',
                    render: function (date) {
                        if (!date) {
                            return '';
                        }

                        return dateToString(date);
                    }
                },
                {
                    data: 'fire_last',
                    render: function (date) {
                        if (!date) {
                            return '';
                        }

                        return dateToString(date);
                    }
                },
                {
                    data: 'events',
                    render: function (data) {
                        var content = '';

                        for (var eventId in data) {
                            content += '<span class="label label-info">' + data[eventId] + '</span> ';
                        }

                        return content;
                    }
                }
            ]
        });

    $('#devices table tbody').on('click', 'button.device-check', function (e) {
        e.preventDefault();
        var deviceId = tableDevices.row($(this).closest('tr')).data().register_id;

        $.ajax({
            type: 'POST',
            url: '/boggart/devices/' + deviceId + '/check',
            success: function() {
                tableDevices.ajax.reload();
            }
        });
    });

    $('#devices table tbody').on('click', 'button.device-ping', function (e) {
        e.preventDefault();
        var device = tableDevices.row($(this).closest('tr')).data();

        $.ajax({
            type: 'POST',
            url: '/boggart/devices/' + device.register_id + '/ping',
            success: function(r) {
                if (r.data !== 'undefined') {
                    if (!r.data) {
                        new PNotify({
                            title: 'Offline',
                            text: 'Device ' + device.id + ' is offline',
                            type: 'error',
                            styling: 'bootstrap3'
                        });
                    } else {
                        new PNotify({
                            title: 'Online',
                            text: 'Device ' + device.id + ' is online',
                            type: 'success',
                            styling: 'bootstrap3'
                        });
                    }
                }
            }
        });
    });

    window.deviceToggle = function(deviceId) {
        $.ajax({
            type: 'POST',
            url: '/boggart/devices/' + deviceId + '/toggle',
            success: function(r) {
                if (r.result === 'failed') {
                    new PNotify({
                        title: 'Error',
                        text: r.message,
                        type: 'error',
                        hide: false,
                        styling: 'bootstrap3'
                    });
                    return
                }

                tableDevices.ajax.reload();
            }
        });
    }
});
$(document).ready(function () {
    var tableDevices = $('#devices table')
        .DataTable({
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
                    data: 'id'
                },
                {
                    data: 'description'
                },
                {
                    data: 'tasks_count'
                },
                {
                    data: 'enabled',
                    render: function (data, type, row) {
                        var content;

                        if (row.enabled) {
                            content = '<button type="button" class="btn btn-danger btn-icon" data-toggle="modal" data-target="#modal" data-modal-title="Confirm disable device #' + row.id + '" data-modal-callback="deviceToggle(\'' + row.register_id + '\');">' +
                                '<i class="glyphicon glyphicon-remove" title="Disable device"></i>'
                        } else {
                            content = '<button type="button" class="btn btn-success btn-icon" data-toggle="modal" data-target="#modal" data-modal-title="Confirm enable device #' + row.id + '" data-modal-callback="deviceToggle(\'' + row.register_id + '\');">' +
                                '<i class="glyphicon glyphicon-ok" title="Enable device"></i>'
                        }

                        return '<div class="btn-group btn-group-xs">' + content + '</button></div>';
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
                    data: 'id'
                },
                {
                    data: 'name'
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
                }
            ]
        });

    window.deviceToggle = function(deviceId) {
        $.ajax({
            type: 'POST',
            url: '/boggart/devices/' + deviceId + '/toggle',
            success: function() {
                tableDevices.ajax.reload();
            }
        });
    }
});
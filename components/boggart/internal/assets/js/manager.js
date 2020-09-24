$(document).ready(function () {
    var groupColumn = 0;

    var tableDevices = $('#devices table')
        .DataTable({
            pageLength: 50,
            language: {
                url: '/dashboard/datatables/i18n.json'
            },
            ajax: {
                url: '/boggart/manager/?entity=devices',
                dataSrc: 'data'
            },
            order: [[groupColumn, 'asc']],
            columnDefs: [{
                visible: false,
                targets: groupColumn
            }],
            columns: [
                {
                    data: 'type'
                },
                {
                    data: 'tags',
                    render: function (tags) {
                        var content = '';

                        for (var i in tags) {
                            content += '<span class="label label-success">' + tags[i] + '</span> ';
                        }

                        return content;
                    }
                },
                {
                    data: 'status',
                    render: function (status) {
                        switch (status.toLowerCase()) {
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
                    data: null,
                    render: function (data, type, row) {
                        var content = '<div class="btn-group" role="group" style="min-width:260px">';

                        if (row.has_widget) {
                            content += '<a href="/boggart/widget/' + row.id + '/" target="_blank" class="btn btn-primary btn-icon btn-xs">' +
                                '<i class="fas fa-window-maximize" title="Open widget"></i>' +
                                '</a>';
                        }

                        if (row.tasks && row.tasks.length > 0) {
                            content += '<div class="btn-group">' +
                                '<button type="button" class="btn btn-primary btn-icon btn-xs" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">' +
                                '<i class="fas fa-running" title="Tasks"></i> <span class="badge">' + row.tasks.length + '</span> <span class="caret"></span>' +
                                '</button>' +
                                '<ul class="dropdown-menu">' +
                                '<li><a href="/boggart/bind/' + row.id + '/tasks/" target="_blank">Show all</a></li>' +
                                '<li role="separator" class="divider"></li>';

                            for (var i in row.tasks) {
                                content +=  '<li><a href="/boggart/bind/' + row.id + '/tasks/?run=' + row.tasks[i][0] + '" target="_blank">Run ' + row.tasks[i][1] + '</a></li>';
                            }

                            content += '</ul></div>';
                        } else if (row.has_readiness_probe || row.has_liveness_probe) {
                            var l = 0;
                            var menu = '';

                            if (row.has_liveness_probe) {
                                l++;
                                menu += '<li><a href="/boggart/bind/' + row.id + '/liveness/" target="_blank">Run liveness probe</a></li>';
                            }

                            if (row.has_readiness_probe) {
                                l++;
                                menu += '<li><a href="/boggart/bind/' + row.id + '/readiness/" target="_blank">Run readiness probe</a></li>';
                            }

                            content += '<div class="btn-group">' +
                                '<button type="button" class="btn btn-primary btn-icon btn-xs" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">' +
                                '<i class="fas fa-running" title="Probes"></i> <span class="badge">' + l + '</span> <span class="caret"></span>' +
                                '</button>' +
                                '<ul class="dropdown-menu">' + menu + '</ul></div>';
                        }

                        if (row.logs_count > 0) {
                            var cl = 'success';

                            switch(row.logs_max_level) {
                                case 'debug':
                                    cl = 'primary';
                                    break;
                                case 'info':
                                    cl = 'info';
                                    break;
                                case 'warn':
                                    cl = 'warning';
                                    break;
                                case 'error':
                                    cl = 'danger';
                                    break;
                                case 'panic':
                                    cl = 'danger';
                                    break;
                                case 'fatal':
                                    cl = 'danger';
                                    break;
                            }

                            content += '<a href="/boggart/bind/' + row.id + '/logs/" target="_blank" class="btn btn-' + cl +' btn-icon btn-xs">' +
                                '<i class="fas fa-headset" title="Show last logs"></i> <span class="badge">' + row.logs_count + '</span>' +
                                '</a>';
                        }

                        if (row.mqtt_publishes > 0 || row.mqtt_subscribers > 0) {
                            content += '<a href="/boggart/bind/' + row.id + '/mqtt/" target="_blank" class="btn btn-primary btn-icon btn-xs">' +
                                '<i class="fas fa-list" title="Show MQTT cache"></i> <span class="badge">' + row.mqtt_publishes + ' | ' + row.mqtt_subscribers + '</span>' +
                                '</a>';
                        }

                        content += '<div class="btn-group">' +
                            '<button type="button" class="btn btn-success btn-icon btn-xs" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">' +
                            '<i class="fas fa-cog" title="Config"></i> <span class="caret"></span>' +
                            '</button>' +
                            '<ul class="dropdown-menu">' +
                            '<li><a href="javascript:void(0)" data-toggle="modal" data-target="#modal" data-modal-title="Device config #' + row.id + '" data-modal-url="/boggart/config/modal/' + row.id + '">Show config</a></li>';

                        if (row.type !== "boggart") {
                            content += '<li><a href="/boggart/bind/' + row.id + '/" target="_blank">Edit bind</a></li>' +
                                '<li><a href="javascript:void(0)" onclick="reloadConfig(\'' + row.id + '\');">Reload from config file</a></li>';
                        }

                        content += '</ul>' +
                            '</div>';

                        if (row.type !== "boggart") {
                            content +=
                                '<button type="button" class="btn btn-danger btn-icon btn-xs" data-toggle="modal" data-target="#modal" data-modal-title="Confirm unregister device #' + row.id + '" data-modal-callback="bindUnregister(\'' + row.id + '\');">' +
                                    '<i class="fas fa-trash" title="Unregister bind"></i>' +
                                '</button>';
                        }

                        return content + '</div>'
                    }
                },
                {
                    data: 'mac'
                },
                {
                    data: 'serial_number'
                },
                {
                    data: 'id'
                },
                {
                    data: 'description'
                }
            ],
            'drawCallback': function () {
                var api = this.api();
                var rows = api.rows({page: 'current'}).nodes();
                var last = null;

                api.column(groupColumn, {page: 'current'}).data().each(function (group, i) {
                    var $aRow = $(rows).eq(i);

                    if (last !== group) {
                        $aRow.before('<tr class="group"></td><td colspan="' + $aRow.children().length + '">' + group + '</td></tr>');
                        last = group;
                    }
                });
            }
        });
    tableDevices.on('click', 'tr.group', function () {
        var currentOrder = tableDevices.order()[0];

        if (currentOrder[0] === groupColumn && currentOrder[1] === 'asc') {
            tableDevices.order([groupColumn, 'desc']).draw();
        }
        else {
            tableDevices.order([groupColumn, 'asc']).draw();
        }
    });

    window.bindUnregister = function (id) {
        $.ajax({
            type: 'POST',
            url: '/boggart/bind/' + id + '/unregister/',
            success: function (r) {
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

                refreshTables();
            }
        });
    };

    window.reloadConfig = function (id) {
        var url = '/boggart/config/reload';

        if (typeof id !== 'undefined') {
            url += '/' + id
        }

        $.ajax({
            type: 'POST',
            url: url,
            success: function (r) {
                if (r.result === 'failed') {
                    new PNotify({
                        title: 'Error',
                        text: r.message,
                        type: 'error',
                        hide: false,
                        styling: 'bootstrap3'
                    });
                } else if (r.message !== 'undefined') {
                    new PNotify({
                        title: 'Success',
                        text: r.message,
                        type: 'success',
                        hide: false,
                        styling: 'bootstrap3'
                    });
                }

                refreshTables();
            }
        });
    };

    window.refreshTables = function () {
        tableDevices.ajax.reload();
    };
});

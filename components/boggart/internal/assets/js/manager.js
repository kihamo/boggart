$(document).ready(function () {
    var groupColumn = 0;

    var tableDevices = $('#devices table')
        .DataTable({
            stateSave: true,
            stateDuration: 0,
            pageLength: 50,
            language: {
                url: '/dashboard/datatables/i18n.json?locale=' + window.shadowLocale
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
                        var cl = 'default';

                        switch (status.toLowerCase()) {
                            case 'online':        cl = 'success'; break;
                            case 'offline':       cl = 'warning'; break;
                            case 'removing':      cl = 'warning'; break;
                            case 'removed':       cl = 'danger';  break;
                            case 'initializing':  cl = 'warning'; break;
                            case 'uninitialized': cl = 'danger';  break;
                        }

                        return '<span class="label label-' + cl + '">' + status + '</span>';
                    }
                },
                {
                    data: null,
                    render: function (data, type, row) {
                        var content = '<div class="btn-group" role="group" style="min-width:380px">';

                        if (row.has_widget) {
                            content += '<a href="/boggart/widget/' + row.id + '/" target="_blank" class="btn btn-primary btn-icon btn-xs">' +
                                '<i class="fas fa-window-maximize" title="Open widget"></i>' +
                                '</a>';
                        }

                        if (row.link.length > 0) {
                            content += '<a href="' + row.link + '" target="_blank" class="btn btn-info btn-icon btn-xs">' +
                                '<i class="fas fa-globe" title="Open link"></i>' +
                                '</a>';
                        }

                        if (row.tasks && row.tasks.length > 0) {
                            var
                                unregistered = 0,
                                listContent = '';

                            for (var i in row.tasks) {
                                if (row.tasks[i].registered) {
                                    listContent += '<li><a href="/boggart/workers/?id=' + row.tasks[i].id + '&action=run" target="_blank">' + row.tasks[i].name;
                                } else {
                                    listContent += '<li><a href="javascript:void(0)"><del>' + row.tasks[i].name + '</del>';
                                    unregistered++;
                                }

                                if (row.tasks[i].custom_schedule) {
                                    listContent += ' <span class="badge">schedule</span>';
                                }

                                listContent += '</a></li>';
                            }

                            content += '<div class="btn-group">' +
                                '<button type="button" class="btn btn-primary btn-icon btn-xs" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">' +
                                '<i class="fas fa-running" title="Tasks"></i> <span class="badge">' + row.tasks.length + ' | ' + unregistered + '</span> <span class="caret"></span>' +
                                '</button>' +
                                '<ul class="dropdown-menu">' +
                                '<li><a href="/boggart/workers/?search=bind/' + row.type + '/' + row.id + '/" target="_blank">Show all</a></li>' +
                                '<li role="separator" class="divider"></li>' + listContent + '</ul>' +
                                '</div>';
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

                            content += '<a href="/boggart/logs/' + row.id + '/" target="_blank" class="btn btn-' + cl +' btn-icon btn-xs">' +
                                '<i class="fas fa-headset" title="Show last logs"></i> <span class="badge">' + row.logs_count + '</span>' +
                                '</a>';
                        }

                        if (row.has_metrics) {
                            content += '<a href="/boggart/metrics/' + row.id + '/" target="_blank" class="btn btn-icon btn-xs';

                            if (row.metrics_empty_count > 0) {
                                content += ' btn-warning';
                            } else {
                                content += ' btn-primary';
                            }

                            content += '">' +
                                '<i class="fas fa-thermometer-empty" title="Show metrics"></i> <span class="badge">' + row.metrics_descriptions_count + ' | ' + row.metrics_collect_count + ' | ' + row.metrics_empty_count + '</span>' +
                                '</a>';
                        }

                        if (row.mqtt_publishes > 0 || row.mqtt_subscribers > 0) {
                            content += '<a href="/boggart/mqtt/' + row.id + '/" target="_blank" class="btn btn-primary btn-icon btn-xs">' +
                                '<i class="fas fa-list" title="Show MQTT cache"></i> <span class="badge">' + row.mqtt_publishes + ' | ' + row.mqtt_subscribers + '</span>' +
                                '</a>';
                        }

                        content += '<div class="btn-group">' +
                            '<button type="button" class="btn btn-success btn-icon btn-xs" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">' +
                            '<i class="fas fa-cog" title="Config"></i> <span class="caret"></span>' +
                            '</button>' +
                            '<ul class="dropdown-menu">' +
                            '<li><a href="javascript:void(0)" data-toggle="modal" data-target="#modal" data-modal-title="Device config #' + row.id + '" data-modal-url="/boggart/config/modal/' + row.id + '/?short=1">Show config</a></li>';

                        if (row.type !== "boggart") {
                            content += '<li><a href="/boggart/bind/' + row.id + '/" target="_blank">Edit bind</a></li>' +
                                '<li><a href="javascript:void(0)" onclick="reloadConfig(\'' + row.id + '\');">Reload from config file</a></li>';
                        }

                        if (row.installers.length > 0) {
                            for (var i in row.installers) {
                                content += '<li><a href="/boggart/installer/' + row.id + '/' + row.installers[i] + '/" target="_blank">Installer ' +
                                    row.installers[i].charAt(0).toUpperCase() + row.installers[i].slice(1) + '</a></li>';
                            }
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
            drawCallback: function () {
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
            },
            initComplete: function() {
                var vars = {};
                window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
                    vars[key] = value;
                });

                if (!vars.hasOwnProperty('search')) {
                    vars['search'] = '';
                }

                this.api().search(vars['search']).draw();
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

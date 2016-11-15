<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link type="text/css" rel="stylesheet" href="/stylesheets/bootstrap.min.css"/>
    <link type="text/css" rel="stylesheet" href="/stylesheets/md.css"/>
    <link type="text/css" rel="stylesheet" href="/stylesheets/metisMenu.min.css"/>
    <link type="text/css" rel="stylesheet" href="/stylesheets/doc.css"/>
</head>

<body>
<div class="docx-wrapper">
    <div class="docx-body clearfix">
        <nav class="navbar navbar-inverse navbar-fixed-top">
            <div class="container-fluid">
                <div class="navbar-header">
                    <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#sidebar-collapse">
                        <span class="sr-only">Toggle navigation</span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                    </button>
                </div>
            </div><!-- /.container-fluid -->
        </nav>

        <div id="sidebar-collapse" class="col-sm-3 col-lg-3 sidebar" role="navigation">
            <div>
                <form role="search" action="javascript:void(0);">
                    <div class="form-group docx-searchForm">
                        <input type="text" class="form-control docx-searchkey" placeholder="Search" name="key" autocomplete="off" >
                        <div class="docx-sug">
                            <ul class="docx-sugul">
                                <li><a href="#">1</a></li>
                                <li><a href="#">1</a></li>
                                <li><a href="#">1</a></li>
                            </ul>
                        </div>
                    </div>

                </form>
            </div>
        </div><!--/.sidebar-->

        <div class="col-sm-9 col-sm-offset-3 col-lg-9 col-lg-offset-3 main docx-marked-wrap container-fluid">
            <div>
                <div class="row">
                </div><!--/.row-->

                <div class="row">
                    <div class="col-lg-12 docx-fade">
                        <div class="docx-panel docx-panel-default">
                            <div class="markdown-body">
                                <div class="docx-marked">
                                    {{$.mdData}}
                                </div>
                            </div>
                        </div>
                    </div>
                </div><!--/.row-->
            </div>
        </div>	<!--/.main-->
    </div>
</div>
<script src="/javascripts/jq.min.js"></script>
<script src="/javascripts/metisMenu.min.js"></script>
<script src="/javascripts/jq.pjax.js"></script>
<script src="/javascripts/bootstrap.min.js"></script>
<script src="/javascripts/doc.js"></script>
</body>

</html>
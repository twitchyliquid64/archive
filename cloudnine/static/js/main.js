google.load("visualization", "1", {packages:["corechart"]});
    function searchSuburb(suburb)
    {
        $("#mostlike").html($('#input-search').val());
        $('#input-search').val(suburb);
        updatePage(true);
        $('#explore').click();
    };

    function compareSuburb(suburb)
    {

    }

$(function(){

    /*-------------------------------------------------------------------*/
    /*  1. Preloader. Requires jQuery jpreloader plugin.
    /*-------------------------------------------------------------------*/
    $(document).ready(function() {
        $('body').jpreLoader({
            showPercentage: false,
            loaderVPos: '50%'
        });
        
        google.setOnLoadCallback(drawPieChart);
        //google.setOnLoadCallback(drawBarChart);

        feed.run();

        $.ajax({
            url: '/api/locations',
            context: document.body
        }).done(function(data){
            locations = JSON.parse(data);
            $( "#input-search" ).autocomplete({
                source: locations
            });
        });
    });

    function drawPieChart() {
        Piedata = google.visualization.arrayToDataTable([
          ['Population', 'Percentage'],
          ['Children',     60],
          ['Young Adults',      25],
          ['Middle Aged',  15],
          ['Elderly',  15]
        ]);

        var options = {
          title: 'Age Distribution',
          titleTextStyle: {color:'white', fontSize: 18},
          pieHole: 0.4,
          colors: ['#3D4A78', '#c80a48', '#FA7AA4', '#F2A2BC', '#F5BFD1', '#D1B6BF', '#BDB3B6'],
          backgroundColor: '#0d0c0d',
          // legend: { position: 'bottom' ,textStyle: {
          //       color: 'white'
          //   }
          legend: { position: 'none' }
        }
        
        lolopt = options;
        var chart = new google.visualization.PieChart(document.getElementById('donutchart'));
        pieC = chart;
        chart.draw(Piedata, options);
    };


    function drawBarChart() {
        var data = google.visualization.arrayToDataTable([
          ['Rating', 'Percentage'],
          ['Suburb',     84],
          ['NSW',     80],
        ]);

        var options = {
            title: 'Health Rating',
            titleTextStyle: {color:'white', fontSize: 18},
            colors: ['#3D4A78'],
            textStyle: {
                color: 'white'
            },
            legend: {position:'none'},
            backgroundColor: '#0d0c0d',
            textStyle: {color: 'white'},
            hAxis: {
                textStyle:{color: '#FFF'}
            },
            vAxis: {
                textStyle:{color: '#FFF'}
            }
        };

        var chart = new google.visualization.BarChart(document.getElementById('barchart'));
        chart.draw(data, options);
    }
    
    
    /*-------------------------------------------------------------------*/
    /*  2. Makes the height of all selected elements (".match-height")
    /*  exactly equal. Requires jQuery matchHeight plugin.
    /*-------------------------------------------------------------------*/
    $(window).smartload(function(){
        $('.match-height').matchHeight();
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  3. Page scrolling feature, requires jQuery Easing plugin.
    /*-------------------------------------------------------------------*/
    var pageScroll = function(){
        $('.page-scroll a').bind('click', function(e){
            e.preventDefault();

            var $anchor = $(this);
            var offset = $('body').attr('data-offset');

            $('html, body').stop().animate({
                scrollTop: $($anchor.attr('href')).offset().top - (offset - 1)
            }, 1500, 'easeInOutExpo');
        });
    };
    
    pageScroll();
    
    
    /*-------------------------------------------------------------------*/
    /*  4. Make navigation menu on your page always stay visible.
    /*  Requires jQuery Sticky plugin.
    /*-------------------------------------------------------------------*/
    var stickyMenu = function(){
        var nav = $('.navbar.navbar-fixed-top');
        nav.unstick();
        nav.sticky({topSpacing: 0});
    };
    
    stickyMenu();
    
    // Call pageScroll() and stickyMenu() when window is resized.
    $(window).smartresize(function(){
        pageScroll();
        stickyMenu();
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  5. Portfolio gallery. Requires jQuery Magnific Popup plugin.
    /*-------------------------------------------------------------------*/
    $('.portfolio').magnificPopup({
        delegate: 'a.zoom',
        type: 'image',
        fixedContentPos: false,
        
        // Delay in milliseconds before popup is removed
        removalDelay: 300,
        
        // Class that is added to popup wrapper and background
        mainClass: 'mfp-fade',
        
        gallery: {
            enabled: true,
            preload: [0,2],
            arrowMarkup: '<button title="%title%" type="button" class="mfp-arrow mfp-arrow-%dir%"></button>',
            tPrev: 'Previous Project',
            tNext: 'Next Project'
        }
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  6. Column Chart (Section - My Strenghts)
    /*-------------------------------------------------------------------*/
    var columnChart = function (){
        $('.column-chart > .chart > .item > .bar > .item-progress').each(function(){
            var item = $(this);
            var newHeight = $(this).parent().height() * ($(this).data('percent') / 100);
            
            // Only animate elements when using non-mobile devices    
            if (jQuery.browser.mobile === false){
                $('.column-chart').one('inview', function(isInView) {
                    if (isInView){
                        // Animate item
                        item.animate({
                            height: newHeight
                        },1500);
                    }
                });
            }
            else{
                item.css('height', newHeight);
            }
        });
    };
    
    // Call columnChart() when window is loaded.
    $(window).smartload(function(){
        columnChart();
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  7. Section - My Resume
    /*-------------------------------------------------------------------*/
    var resumeCollapse = function (){
        var workItem = $('#work .collapse:not(:first)');
        var educationItem = $('#education .collapse:not(:first)');
        var ww = Math.max($(window).width(), window.innerWidth);

        if (ww < 768){
            workItem.collapse('show');
            educationItem.collapse('show');
        }
        else{
            workItem.collapse('hide');
            educationItem.collapse('hide');
        }
    };
    
    // Call resumeCollapse() when window is loaded.
    $(window).smartload(function(){
        resumeCollapse();
    });
    
    // Call resumeCollapse() when window is resized.
    $(window).smartresize(function(){
        resumeCollapse();
    });
    
    
    /*-------------------------------------------------------------------*/
    /*	8. References slider. Requires Flexslider plugin.
    /*-------------------------------------------------------------------*/
    $(window).smartload(function(){
        var flex = $('.flexslider.references');
    
        flex.flexslider({
            selector: ".slides > .item",
            manualControls: ".flex-control-nav li",
            directionNav : false,
            slideshowSpeed: 4000,
            after: function(slider){
                if (!slider.playing) {
                    slider.play();
                }
            }
        }); 
    });
    
    $('a.flex-prev').on('click', function(e){
        e.preventDefault();
        $('.flexslider').flexslider('prev');
    });
    
    $('a.flex-next').on('click', function(e){
        e.preventDefault();
        $('.flexslider').flexslider('next');
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  9. Circle Chart (Section - Skills & Expertise)
    /*-------------------------------------------------------------------*/
    var circleChart = function (){
        $('.circle-chart .item > .circle > .item-progress').each(function(){
            var item = $(this);
            var maxHeight = 108;
            var newHeight = maxHeight * ($(this).data('percent') / 100);
            
            // Only animate elements when using non-mobile devices    
            if (jQuery.browser.mobile === false){
                item.one('inview', function(isInView) {
                    if (isInView){
                        // Animate item
                        item.animate({
                            height: newHeight
                        },1500);
                    }
                });
            }
            else{
                item.css('height', newHeight);
            }
        });
    };
    
    // Call circleChart() when window is loaded.
    $(window).smartload(function(){
        circleChart();
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  10. Bar Chart (Section - Knowledge)
    /*-------------------------------------------------------------------*/
    var barChart = function (){
        $('.bar-chart > .item > .bar > .item-progress').each(function(){
            var item = $(this);
            var percent = $(this).prev();
            var newWidth = $(this).parent().width() * ($(this).data('percent') / 100);
            
            // Only animate elements when using non-mobile devices    
            if (jQuery.browser.mobile === false){
                item.one('inview', function(isInView) {
                    if (isInView){
                        // Animate item
                        item.animate({
                            width: newWidth
                        },1500);
                        
                        percent.animate({
                            left: newWidth - percent.width()
                        },1500);
                    }
                });
            }
            else{
                item.css('width', newWidth);
                percent.css('left', newWidth - percent.width());
            }
        });
    };
    
    // Call barChart() when window is loaded.
    $(window).smartload(function(){
        barChart();
    });
    
    // Call barChart() when window is resized.
    $(window).smartresize(function(){
        barChart();
    });
    
    
    /*-------------------------------------------------------------------*/
    /*  11. Milestones counter.
    /*-------------------------------------------------------------------*/
    var counter = function (){
        var number = $('.milestones .number');
        
        number.countTo({
            speed: 3000
        });
    };
    
    if (jQuery.browser.mobile === false){
        var number = $('.milestones .number');
        
        number.one('inview', function(isInView) {
            if (isInView){
                counter();
            }
        });
    }
    else{
        counter();
    }

    var feed = new Instafeed({
        get: 'tagged',
        tagName: 'OperaHouse',
        clientId: '36b861cbb3da477d926f193f3e324568',
        limit: 6
    });

    /*-------------------------------------------------------------------*/
    // API for our system.
    /*-------------------------------------------------------------------*/
    locations = [];


    
    updatePage = function(forcePage)
    {
        var LGA = $('#input-search').val();

	var overrideSelection = forcePage || false;
	if (overrideSelection)
		URL = '/api/display?lga=' + LGA;
	else
		URL = '/api/similar?lga=' + LGA

	$("#loadAni").css("visibility","initial");
        console.log(LGA);
        $.ajax({
            url: URL,
            context: document.body
        }).done(function(data){
    	    $("#loadAni").css("visibility","hidden");
                d = JSON.parse(data);
            
	    if(!overrideSelection)
    	        $("#mostlike").html($('#input-search').val());

	    $("#locspot").html(d[0].name);
    	    $("#otherlist").html("");
    	    for( var i = 1; i < d.length; i++)
                $("#otherlist").append(
                    // $('<tr><td>'+d[i].name+'</td>   <td><a href="#profile" onclick="searchSuburb(\''+d[i].name+'\');"><span class="icons icons-explore" data-toggle="tooltip" data-placement="top" title="Explore this suburb."></span></a></td>   <td><a href="#services" onclick="searchSuburb(\''+d[i].name+'\');"><span class="icons icons-compare"  data-toggle="tooltip" data-placement="top" title="Compare with this suburb."></span></a></td></tr>'));
                    $('<tr><td>'+d[i].name+'</td>   <td><a href="#profile" onclick="searchSuburb(\''+d[i].name+'\');">Explore</a></td>   <td><a href="#services" style="visibility:hidden" onclick="compareSuburb(\''+d[i].name+'\');">Compare</a></td></tr>'));
    	           
            Piedata.setCell(0,1,d[0].population[0]);
            Piedata.setCell(1,1,d[0].population[1]);
            Piedata.setCell(2,1,d[0].population[2]);
            Piedata.setCell(3,1,d[0].population[3]);
            pieC.draw(Piedata, lolopt);


	    if (d[0].rank)
		    $("#rank").text(d[0].rank[1].formatMoney(0,',',','));
	    else
		    $("#rank").text("N/A");
	    $("#population").text((d[0].population[0] + d[0].population[1] + d[0].population[2] + d[0].population[3]).formatMoney(0,',',','));
    	    var c = d[0];
    	    var assault = 0;
    	    var drug = 0;
    	    var theft = 0;
    	    var robbery = 0;
    	    for(var i = 0;i < c.crime.length;i++)
    	    {
    		if (c.crime[i][2] == "Assault"){
    			assault += parseInt(c.crime[i][4]);
    		};
                    if (c.crime[i][2] == "Drug offences"){
                            drug += parseInt(c.crime[i][4]);
                    };
                    if (c.crime[i][2] == "Robbery"){
                            robbery += parseInt(c.crime[i][4]);
                    };
                    if (c.crime[i][2] == "Theft"){
                            theft += parseInt(c.crime[i][4]);
                    };
    	    }
    	    var max = Math.max(assault,drug,robbery,theft);
    	    // console.log(max,assault,drug,robbery,theft);	
            var countLength = Math.pow(10, max.toString().length - 1);
            var maxRound = Math.ceil(max/countLength)*countLength;
            var assaultPers = Math.round(assault / maxRound * 320);
            var theftPers = Math.round(theft / maxRound * 320);
            var drugPers = Math.round(drug/ maxRound * 320);
            var robberyPers = Math.round(robbery / maxRound * 320);

    	    
    	    // console.log(assaultPers, theftPers);
    	    $("#assault").css("height", assaultPers < 30 ? 30 : assaultPers);
    	    $("#theft").css("height", theftPers < 30 ? 30 : theftPers);
    	    $("#robbery").css("height", drugPers < 30 ? 30 : robberyPers);
    	    $("#drug").css("height", robberyPers < 30 ? 30 : drugPers);
            $('#per-100').html(maxRound);
            $('#per-75').html(maxRound*0.75);
            $('#per-50').html(maxRound*0.5);
            $('#per-25').html(maxRound*0.25);
            $('#per-assault').text(assault);
            $('#per-drug').text(drug);
            $('#per-robbery').text(robbery);
            $('#per-theft').text(theft);

            var word = getRandom6Words(exploreWords);

            $('#cloudtag-1').html(word[0]);
            $('#cloudtag-2').html(word[1]);
            $('#cloudtag-3').html(word[2]);
            $('#cloudtag-4').html(word[3]);
            $('#cloudtag-5').html(word[4]);
            $('#cloudtag-6').html(word[5]);

            var countAccolades = 0;
            $.ajax({
                    url: '/api/stories?lga=' + LGA,
                    context: document.body
            }).done(function(data){
                data = JSON.parse(data);
                $(".accolades").html("");
                countAccolades = data.length;
                for(var i = 0; i < countAccolades; i++)
                {
                    var newd = '<div class="item"><p class="normal-text">' + data[i][2].substring(data[i][2].length - 4) + '</b><div class="content"><a href="' + data[i][0] + '" style="color: #FFF;"><h3>' + data[i][1] + '</h3></a></div></div>';
                    $(".accolades").append($(newd));
                }

                console.log(countAccolades, c.name);
                $('#instafeed').empty();
                feed.options.tagName = c.name.split(" ")[0];
                feed.options.limit = countAccolades > 3 ? 4 : 6;
                feed.run();
            });
            
        });
    };
    function isValidLocation()
    {
        for(var i = 0; i < locations.length; i++)
        {
            if (locations[i] == $('#input-search').val())
                return true;
        }
        return false;
    };

    $('#input-search').keypress(function(event) {
        if (event.keyCode == 13 && isValidLocation()) {
            updatePage(true);
        }
    });      

    $('#button-submit').click(function(){
        if (isValidLocation()) {
            updatePage(true);
        }
    });


    var words = ["crime statistics", "population", "news stories", "health indicators", "lifestyle", "household types", "gender split"
            , "country of birth", "median age", "languages spoken at home", "household income", "education"];

    var exploreWords = ["festival", "community", "scholarship", "sport", "wins", "competition", "family", "anniversary", "opening", "ceremory"
    , "unveil", "marriage", "event", "culture", "education"];

    getRandom6Words = function (awords) {
    	if (awords === undefined) awords = words;
    	
        var result = [];
        var inputedIndex = [];
        for(var i=0; i<6; i++) {
            var position = Math.floor(Math.random() * awords.length);
            while(inputedIndex[position]) {
                position = Math.floor(Math.random() * awords.length);
            }
            inputedIndex[position] = true;
            result[i] = awords[position];
        }
        return result;
    };

    Number.prototype.formatMoney = function(c, d, t){
    var n = this, 
        c = isNaN(c = Math.abs(c)) ? 2 : c, 
        d = d == undefined ? "." : d, 
        t = t == undefined ? "," : t, 
        s = n < 0 ? "-" : "", 
        i = parseInt(n = Math.abs(+n || 0).toFixed(c)) + "", 
        j = (j = i.length) > 3 ? j % 3 : 0;
       return s + (j ? i.substr(0, j) + t : "") + i.substr(j).replace(/(\d{3})(?=\d)/g, "$1" + t) + (c ? d + Math.abs(n - i).toFixed(c).slice(2) : "");
     };
    
    $(document).ready(function() { when_ready(); });
});

$(document).ready(function () {
    var $form = $('#form');
    var $input = $('#url');
    var $pageSum = $('#summary');
    var $error = $('#error');
    let baseURL = "https://api.patrick-yi.com/v1/summary?url=";


    $form.submit(function (e) {
        e.preventDefault();

        // Remove previously displayed page summary if any.
        $pageSum.html('');
        $error.html('');

        var url = $input.val();
        url = encodeURIComponent(url);
        console.log("url" + url)

        // should change
        url = baseURL + url;

        $.ajax({
                url: url
            })
            .done(function (data) {
                var summary = "";
                var images = "";
                var title = "";
                var desc = "";

                for (each in data) {
                    switch (each) {
                        case 'title':
                            title = '<h3>' + data[each] + '</h3>';
                            break;
                        case 'description':
                            desc = '<p>' + data[each] + '</p>';
                            break;
                        case 'images':
                            data[each].forEach(function (img) {
                                if (!img['url']) {
                                    return;
                                }
                                images += '<img src="' + img['url'] + '"' + '>';
                            });
                            break;
                        default:
                            break;
                    }
                }
                summary = title + desc + images;
                $pageSum.html(summary);
            })
            .fail(function (error) {
                $error.text(error.responseText);
            });
    });

});
//  common function
function isYYYYMMDD(stro){
        let str = '';
        str = jQuery.trim(stro);
        if (str == '') {
          return true;
        }
        //null or 8文字でない or 数値でない場合はfalse
        if(str.length != 8 || isNaN(str)){
          return false;
        }
        //年,月,日を取得する
        var y = parseInt(str.substr(0,4));
        var m = parseInt(str.substr(4,2)) -1;  //月は0～11で指定するため-1しています。
        var d = parseInt(str.substr(6,2));
        var dt = new Date(y, m, d);

        //判定する
        return (y == dt.getFullYear() && m == dt.getMonth() && d == dt.getDate());
};


function iscompare(fr, to){
//  console.log( "fr is " + fr);
//  console.log( "to is " + to);
  var d1 = Number(fr);
  var d2 = Number(to);
//  console.log("d1 is " + d1);
//  console.log("d2 is " + d2);
  if (d1 <= d2) {
    return true
  } else {
    return false
  }; // true
}

function file_upload(){
    // フォームデータを取得
    var formdata = new FormData($('#upload_form').get(0));
    console.log(formdata);

    // POSTでアップロード
    $.ajax({
        url  : "/uploads/",
        type : "POST",
        data : formdata,
        cache       : false,
        contentType : false,
        processData : false,
        dataType    : "json"
    })
    .done(function(result, textStatus, xhr){
//        alert(data);

        console.log(result);
//        let res = JSON.parse(result);

        $("#mini-tail2").html('');
        $("#mini-tail2").html('アップロードに成功しました。' + result.results);


    })
    .fail(function(xhr, textStatus, error){
      // 異常系の処理を書く
      $("#mini-tail2").html('');
      $("#mini-tail2").html('アップロードに失敗しました。');
    })
    .always(function(xhr, settings) {
    // 正常・異常系に関わらず行われる処理を書く
    });
}


$(document).ready(function(){

//  at jquery object clicked

    $("#btn-search").on('click', function() {    //  When clicking on select button

        // check process : check values in each of fields.
        if (isYYYYMMDD($('input[name="fromdate"]').val()) != true) {
          alert('Typo error on from date');
          return false
        }
        if (isYYYYMMDD($('input[name="todate"]').val()) != true) {
          alert('Typo error on to-date');
          return false
        }

        if (iscompare($('input[name="fromdate"]').val(), $('input[name="todate"]').val()) == false ) {
          alert('To-date should be equal or greater than from-date.');
          return false
        }


        // data transfer process : POST

        var search_obj = new Object();
        search_obj.callsign = $('input[name="callsign"]').val();
        search_obj.frdate = $('input[name="fromdate"]').val();
        search_obj.todate = $('input[name="todate"]').val();
        var search_json = '';
        search_json = JSON.stringify(search_obj);
         console.log(search_json);
        $.ajax({
      		url : '/',
      		type : 'POST',
      		data : search_json,
          dataType: 'json',
      		async : true,
      		timeout : 10000,

          // data reception process

      	}).done(function(result, textStatus, xhr) {
      		// 正常系の処理を書く
//          console.log(result);


          var bodyhead = $('<div id="tablebefore">以下に検索した結果を示します。ひとつのQSOに対して２つのpdfファイルになることもあります。</div>')
          var tabledef = $('<table class="table" id="down-table" bordercolor="#000000">');
          var thead = $('<thead color="blue">');
          var tbody = $('<tbody id="down-tbody">');
          var bodytail = $('<div id="preparing-file-modal" title="Preparing report..." style="display: none;"> We are preparing your report, please wait... <div class="ui-progressbar-value ui-corner-left ui-corner-right" style="width: 100%; height:22px; margin-top: 20px;"></div></div><div id="error-modal" title="Error" style="display: none;">There was a problem generating your report, please try again.</div>');
          var tr_text = '';
          tabledef.addClass('table-striped'); //Bootstrapのテーブルを適用するときは必ずtbodyタグを使用する）
          thead.append('<tr><th>コールサイン</th><th>日付</th><th>ダウンロード用リンクボタン</th></tr>');

          console.log(result.results);
          let res = JSON.parse(result.results);
          for (let item in res) {
              console.log(res[item]);
              wk_date = new String(res[item].Date);
              tr_text = ''
              tr_text = '<tr>';
              tr_text = tr_text +'<td>' + res[item].Callsign + '</td>';
              tr_text = tr_text + '<td>' + wk_date.substr(0,11) + '</td>';
              tr_text = tr_text + '<td>' + '<a class="fileDownloadCustomRichExperience" href="/downloads/'+ res[item].File +'">'+ res[item].File + '</a>' + '</td>';
              tr_text = tr_text + '</tr>';
              tbody.append(tr_text);
          }

          tabledef.append(thead);
          tabledef.append(tbody);
          $("#mini-tail").html('');
          $("#mini-tail").append(bodyhead);
          $("#mini-tail").append(tabledef);
          $("#mini-tail").append(bodytail);

      	}).fail(function(xhr, textStatus, error) {
      		// 異常系の処理を書く
          $("#mini-tail").html('');
      	}).always(function(xhr, settings) {
      		// 正常・異常系に関わらず行われる処理を書く
      	});
      });   // end of btn-search

      $(document).on("click", "a.fileDownloadCustomRichExperience", function () {

//           var $preparingFileModal = $("#preparing-file-modal");

//           $preparingFileModal.dialog({ modal: true });

           $.fileDownload($(this).prop('href'), {
               successCallback: function (url) {

  //                 $preparingFileModal.dialog('close');
  //              alert('You just got a file download dialog or ribbon for this URL :' + url);
               },
               failCallback: function (responseHtml, url) {

//                 alert('Your file download just failed for this URL:' + url + '\r\n' +
//                         'Here was the resulting error HTML: \r\n' + responsehtml
//                         );
//                   $preparingFileModal.dialog('close');
                   $("#error-modal").dialog({ modal: true });
               }
           });
           return false; //this is critical to stop the click event which will trigger a normal file download!
       });
});   //  End of document ready

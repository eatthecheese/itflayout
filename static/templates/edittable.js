function edit_row(id) {
    var cur_location=document.getElementById("location"+id).innerHTML;

    document.getElementById("location"+id).innerHTML="<input type='text' id='location_text"+id+"' value='"+cur_location+"'>";

    document.getElementById("edit_button"+id).style.display="none";
}
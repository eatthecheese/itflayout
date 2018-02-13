function edit_row(id) {
    var cur_ip=document.getElementById("ip"+id).innerHTML;
    var cur_location=document.getElementById("location"+id).innerHTML;
    var cur_version=document.getElementById("version"+id).innerHTML;
    var cur_nlc=document.getElementById("nlc"+id).innerHTML;
    var cur_scnumber=document.getElementById("scnumber"+id).innerHTML;
    var cur_transportmode=document.getElementById("transportmode"+id).innerHTML;
    var cur_environment=document.getElementById("environment"+id).innerHTML;
    var cur_priconc=document.getElementById("priconc"+id).innerHTML;
    var cur_secconc=document.getElementById("secconc"+id).innerHTML;
    //var cur_devicesactive=document.getElementById("devicesactive"+id).innerHTML;
    var cur_id = id.toString();

    document.getElementById("ip"+id).innerHTML="<input type='text' id='ip_text"+id+"' name='new_ip"+id+"' value='"+cur_ip+"'>";
    document.getElementById("formdummy").innerHTML="<input type='hidden' name='id' value='"+cur_id+"'>";
    document.getElementById("location"+id).innerHTML="<input type='text' id='location_text"+id+"' name='new_location"+id+"' value='"+cur_location+"'>";
    document.getElementById("version"+id).innerHTML="<input type='text' id='version_text"+id+"' name='new_version"+id+"' value='"+cur_version+"'>";
    document.getElementById("nlc"+id).innerHTML="<input type='text'id='nlc_text"+id+"' name='new_nlc"+id+"' value='"+cur_nlc+"'>";
    document.getElementById("scnumber"+id).innerHTML="<input type='text' id='scnumber_text"+id+"' name='new_scnumber"+id+"' value='"+cur_scnumber+"'>";
    document.getElementById("transportmode"+id).innerHTML="<input type='text' id='transportmode_text"+id+"' name='new_transportmode"+id+"' value='"+cur_transportmode+"'>";
    document.getElementById("environment"+id).innerHTML="<input type='text' id='environment_text"+id+"' name='new_environment"+id+"' value='"+cur_environment+"'>";
    document.getElementById("priconc"+id).innerHTML="<input type='text' id='priconc_text"+id+"' name='new_priconc"+id+"' value='"+cur_priconc+"'>";
    document.getElementById("secconc"+id).innerHTML="<input type='text' id='secconc_text"+id+"' name='new_secconc"+id+"' value='"+cur_secconc+"'>";
    //document.getElementById("devicesactive"+id).innerHTML="<input type='text' id='devicesactive_text"+id+"' name='new_devicesactive"+id+"' value='"+cur_devicesactive+"'>";

    document.getElementById("edit_button0").style.display="none";
    {{range .}}
        document.getElementById("edit_button"+{{.Scid}}).style.display="none";
    {{end}}
    document.getElementById("save_button"+id).style.display="block";
}

function save_row(id) {
    /*
    var new_ip=document.getElementById("ip"+id).value;
    var new_location=document.getElementById("location"+id).value;
    var new_version=document.getElementById("version"+id).value;
    var new_nlc=document.getElementById("nlc"+id).value;
    var new_scnumber=document.getElementById("scnumber"+id).value;
    var new_transportmode=document.getElementById("transportmode"+id).value;
    var new_environment=document.getElementById("environment"+id).value;
    var new_priconc=document.getElementById("priconc"+id).value;
    var new_secconc=document.getElementById("secconc"+id).value;
    var new_devicesactive=document.getElementById("devicesactive"+id).value;
    */

    document.getElementById("formtable").submit();
}      
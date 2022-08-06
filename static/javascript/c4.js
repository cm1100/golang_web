let notify = function(c){

    const {msg="",title="",footer=null,icon='success'}=c
    Swal.fire({
          icon: icon,
          title: title,
          text: msg,
          footer: footer
            })
  }
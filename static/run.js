/*
Copyright (c) 2012, André Simon
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.


THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL ANDRE SIMON BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

var runSlidePort=12345;
var currentOutputNodeName;

function run(type, codeId, template)
{

     var ws = new WebSocket("ws://localhost:"+runSlidePort+"/ws");
     ws.onopen = function()
     {
	var myMap = {};
	myMap["Type"] = type;
	myMap["Template"] = template;
	myMap["Snippet"] = document.getElementById(codeId).textContent;
	currentOutputNodeName=codeId+"Out";
        ws.send(JSON.stringify(myMap));
        ws.close();
     };

     ws.onmessage = function (evt)
     {
      console.log(evt.data);
      var received_msg = evt.data;
      document.getElementById(currentOutputNodeName).firstChild.nodeValue=received_msg;
     };
     ws.onerror   = function (evt) { console.log('Error occured: ' + evt.data); };
     ws.onclose   = function (evt) { console.log("Disconnected"); };

}



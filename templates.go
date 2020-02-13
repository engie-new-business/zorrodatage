package main

import "html/template"

var homeTemplate = template.Must(template.New("home").Parse(homePage))

var homePage = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8" name="viewport" content="width=device-width, initial-scale=1">
    <title>Timestamping</title>
    <link rel="stylesheet" href="https://unpkg.com/purecss@1.0.1/build/pure-min.css" integrity="sha384-oAOxQR6DkCoMliIh8yFnu25d7Eq/PHS21PClpwjOTeU2jRSq11vu66rf90/cZr47" crossorigin="anonymous">
</head>
<style>
    header {
        grid-area: header;
        margin-bottom: 3em;
        text-align: center;
    }
    nav {
        grid-area: nav;
    }
    main {
        grid-area: main;
    }
    footer {
        grid-area: footer;
        text-align: center;
    }
    ul {
        margin: 0;
    }

    .grid-container {
        display: grid;
        grid-template-areas:
                'header header'
                'nav main'
                'footer footer';
        grid-column-gap: 2em;
        grid-template-rows: auto 1fr auto;
        grid-template-columns: auto 1fr;
    }
</style>
<body>
    <section class="grid-container">
        <header>
            <h1>Rockside Timestamping Showcase</h1>
        </header>
        <nav>
            <div class="pure-menu custom-restricted-width">
                <ul class="pure-menu-list">
                    <li class="pure-menu-heading">APPLICATION</li>
                    <li class="pure-menu-item"><a target="_blank" href="https://github.com/rocksideio/zorrodatage/blob/master/README.md" class="pure-menu-link">How does it work?</a></li>
                    <li class="pure-menu-item"><a target="_blank" href="https://github.com/rocksideio/zorrodatage" class="pure-menu-link">Github Code</a></li>
                    <li class="pure-menu-item"><a target="_blank" href="https://github.com/rocksideio/zorrodatage/blob/master/contract.sol" class="pure-menu-link">Ethereum Contract</a></li>
                    <li class="pure-menu-item"><a target="_blank" href="https://github.com/rocksideio/rockside-sdk-go/" class="pure-menu-link">Rockside GO SDK</a></li>
                    <li class="pure-menu-heading"></li>
                    <li class="pure-menu-heading"></li>
                    <li class="pure-menu-heading">ABOUT</li>
                    <li class="pure-menu-item"><a target="_blank" href="https://docs.rockside.io/" class="pure-menu-link">Rockside Documentation</a></li>
                    <li class="pure-menu-item"><a target="_blank" href="https://rockside.io" class="pure-menu-link">Rockside Company</a></li>
                </ul>
            </div>
        </nav>
        <main>
            <div>
                <h3>Store a proof of existence in the blockchain</h3>
                <form class="pure-form" name="file-upload" action="/upload" method="post" enctype="multipart/form-data">
                    <fieldset>
                        <legend>Timestamp your file (less than 1M)</legend>
                        <input type="file" id="file-input" required />
                        <input type="submit" class="pure-button pure-button-primary" value="Upload" />
                    </fieldset>
                </form>
                <ul id="registered-result"></ul>
                <p id="upload-error"></p>
            </div>
            <div>
                <h3>Lookup a fingerprint already stored</h3>
                <fieldset>
                    <legend>Enter your fingerprint</legend>
                    <input type="text" name="fingerprint" size="60" placeholder="ex: 09985b7c0618392bfa8caef726c6bfe0fda42f17360a43a62dbcc58c2ee613bf" required />
                    <button type="submit" class="pure-button pure-button-primary" id="fingerprint-button">Check</button>
                </fieldset>
                <ul id="stamp"></ul>
                <p id="check-error"></p>
            </div>
        </main>
        <footer>&copy; Rockside. 2019 Blockchain Studio - All rights reserved.</footer>
    </section>
</body>

<script>
    var button = document.querySelector("#fingerprint-button");
    button.addEventListener("click", function (event) {
        event.preventDefault();
        const params = new URLSearchParams();
        params.append("fingerprint", document.querySelector("input[name='fingerprint']").value);
        fetch("lookup", {
            method: "POST",
            body: params,
        }).then( response => {
            document.querySelector("#check-error").innerHTML = "";
            if(!response.ok) {throw response}
            return response.json();
        }).then( jsonData => {
            var result = document.querySelector("#stamp");
            result.innerHTML = "";
            var hash = document.createElement("li");
            hash.innerHTML = "<li>Date: " + jsonData.Date+ "</li>";
            var fingerprint = document.createElement("li");
            fingerprint.innerHTML = "<li>File Fingerprint: " + jsonData.Fingerprint + "</li>";
            result.appendChild(fingerprint);
            result.appendChild(hash);
        }).catch( err => {
            err.text().then( message => {
                result = document.querySelector("#check-error");
                result.innerHTML = message;
            })
        });
        event.preventDefault();
    });

    var form = document.forms.namedItem("file-upload");
    form.addEventListener("submit", function (event) {
        formData = new FormData();
        var file = document.querySelector("#file-input").files[0];
        formData.append("file", file);
        fetch("/upload", {
            method: "POST",
            body: formData,
        }).then( response => {
            document.querySelector("#upload-error").innerHTML = "";
            if(!response.ok) {throw response}
            return response.json();
        }).then( jsonData => {
            var result = document.querySelector("#registered-result");
            result.innerHTML = "";
            var hash = document.createElement("li");
            hash.innerHTML = "Tx Hash: <a target=\"_blank\" href=" + jsonData.HashURL+ ">" + jsonData.TransactionHash + "</a>";
            var fingerprint = document.createElement("li");
            fingerprint.innerHTML = "File Fingerprint: " + jsonData.Fingerprint;
            result.appendChild(fingerprint);
            result.appendChild(hash);
        }).catch( err => {
            err.text().then( message => {
                result = document.querySelector("#upload-error");
                result.innerHTML = message;
            })
        });
        event.preventDefault();
    });
</script>
</html>`

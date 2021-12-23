const grpc = require("@grpc/grpc-js");
var protoLoader = require("@grpc/proto-loader");
const PROTO_PATH = "hellow.proto";

const options = {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
};

var packageDefinition = protoLoader.loadSync(PROTO_PATH, options);
const routeguide = grpc.loadPackageDefinition(packageDefinition).hellow;
const client = new routeguide.Greeter(
    "192.168.1.185:30000",
    grpc.credentials.createInsecure()
);
// client.sayHello({ name: "you" }, function (err, response) {
//     console.log("Greeting:", response.message);
// });
// var call = client.ListResponse({ messaqe: "test" });
// call.on("data", function (response) {
//     console.log(response);
// });
// call.on("end", function () {
//     //console.log("end");
// });
// call.on("error", function (e) {
//     // An error has occurred and the stream has been closed.
// });
// call.on("status", function (status) {
//     //console.log(status);
// });
let meta = new grpc.Metadata();
meta.add('authorization', 'token');
var call = client.Stream(meta);
call.on("data", (response) => {
    console.log(response.message);
});

var readline = require('readline');
var rl = readline.createInterface(process.stdin, process.stdout);
rl.on('line', function(line) {
    //if (line === "right") rl.close();
    // if (line == "s"){
    //     call.write({typeS:{mtypes:line}})
    // }
    // if (line == "w"){
    //     call.write({typeW:{mtypew:line}})
    // }
    //console.log("message: "+line)
    call.write({typeS:{mtypes:line}})
}).on('close',function(){
    process.exit(0);
});
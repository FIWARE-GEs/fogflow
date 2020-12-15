(function() {

function CtxElement2JSONObject(e) {
    var jsonObj = {};

    for (let key in e) {
	jsonObj[key] = e[key]
    }    

    return jsonObj;
}    

function JSONObject2CtxElement(ctxObj) {
    console.log('convert json object to context element') 

    var ctxElement = {}
    ctxElement['id'] = ctxObj['id'] 
    ctxElement['type'] = ctxObj['type']
    for (let key in ctxObj) {
        if ((key != 'id') && (key != 'type') && (key != 'modifiedAt') && (key != 'createdAt') && (key != 'observationSpace') && (key != 'operationSpace') && (key != 'location') && (key != '@context')) {
	    ctxElement[key] = ctxObj[key]
	}

    }
    return ctxElement
}  

    
var NGSI10Client = (function() {
    // initialized with the broker URL
    var NGSI10Client = function(url) {
	var url = str.substring(0, str.lastIndexOf("/") + 1);
        this.brokerURL = url;
    };
    
    // update context 
    NGSI10Client.prototype.updateContext = function updateContext(ctxObj) {
        updateCtxReq = JSONObject2CtxElement(ctxObj);
	console.log(updateCtxReq);
		      
        return axios({
            method: 'post',
            url: this.brokerURL + '/ngsi-ld/v1/entities/',
            data: updateCtxReq
	    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        }).then( function(response){
            if (response.status == 201) {
                return response.data;
            } else {
                return null;
            }
        });
    };
            
    // subscribe context
    NGSI10Client.prototype.subscribeContext = function subscribeContext(subscribeCtxReq) {        
        return axios({
            method: 'post',
            url: this.brokerURL + '/ngsi-ld/v1/subscriptions/',
            data: subscribeCtxReq
	    headers = {'Accept': 'application/ld+json',
               'Content-Type': 'application/ld+json',
               'Link': '<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        }).then( function(response){
            if (response.status == 201) {
                return response.data.subscribeResponse.subscriptionId;
            } else {
                return null;
            }
        });
    };    


var NGSI9Client = (function() {
    // initialized with the address of IoT Discovery
    var NGSI9Client = function(url) {
        this.discoveryURL = url;
    };
        
    NGSI9Client.prototype.findNearbyIoTBroker = function findNearbyIoTBroker(mylocation, num) 
    {
        var discoveryReq = {};    
        discoveryReq.entities = [{type: 'IoTBroker', isPattern: true}];              
    
        var nearby = {};
        nearby.latitude = mylocation.latitude;
        nearby.longitude = mylocation.longitude;
        nearby.limit = num;
        
        discoveryReq.restriction = {
            scopes: [{
                type: 'nearby',
                value: nearby
            }]
        };
    
        return this.discoverContextAvailability(discoveryReq).then( function(response) {
            if (response.errorCode.code == 200) {
                var brokers = [];
                for(i in response.contextRegistrationResponses) {
                    contextRegistrationResponse = response.contextRegistrationResponses[i];
                    var providerURL = contextRegistrationResponse.contextRegistration.providingApplication;
                    if (providerURL != '') {
                        brokers.push(providerURL);
                    }
                }
                return brokers;
            } else {
                return nil;
            }            
        });
    }
            
    // discover availability
    NGSI9Client.prototype.discoverContextAvailability = function discoverContextAvailability(discoverReq) {        
        return axios({
            method: 'post',
            url: this.discoveryURL + '/discoverContextAvailability',
            data: discoverReq
        }).then( function(response){
            if (response.status == 200) {
                return response.data;
            } else {
                return null;
            }
        });
    };               
    
    return NGSI9Client;
})();

// initialize the exported object for this module, both for nodejs and browsers
if (typeof module !== 'undefined' && typeof module.exports !== 'undefined'){
    this.axios = require('axios')    
    module.exports.NGSI10Client = NGSI10Client; 
    module.exports.NGSI9Client = NGSI9Client;   
    module.exports.CtxElement2JSONObject = CtxElement2JSONObject;
    module.exports.JSONObject2CtxElement = JSONObject2CtxElement;    
} else {
    window.NGSI10Client = NGSI10Client;  
    window.NGSI9Client = NGSI9Client;
    window.CtxElement2JSONObject = CtxElement2JSONObject;
    window.JSONObject2CtxElement = JSONObject2CtxElement;
}

})();
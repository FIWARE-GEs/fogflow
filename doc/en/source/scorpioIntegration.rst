*****************************************
Integrate FogFlow with Other NGSI-LD Broker
*****************************************


Scenario
===============================================


.. figure:: figures/fogflow-scorpio.png


How to set up FogFlow and Scorpio
===============================================

Prerequisite for all System
------------------------------------------------

Here are the prerequisite commands for starting all system :

1. docker

2. docker-compose

For ubuntu system, you need to install docker-ce and docker-compose.

To install Docker CE, please refer to `Install Docker CE`_, required version > 18.03.1-ce;

.. important:: 
	**please also allow your user to execute the Docker Command without Sudo**


To install Docker Compose, please refer to `Install Docker Compose`_, 
required version 18.03.1-ce, required version > 2.4.2

.. _`Install Docker CE`: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04
.. _`Install Docker Compose`: https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04

Setup FogFlow System
------------------------------------------------
* To install FogFlow System, please refer  `FogFlow installation Document`_, 

Setup Other NGSI-LD Broker
------------------------------------------------
* To install Orion Broker, please refer to `Docker Compose file for Orion-Broker`_, 
* To install Scorpio Broker, please refer to `Docker Compose file for Scorpio Broker`_, 
* To install  stellio-context-broker-ld, please refer to `Docker Compose file for Stellio-context-broker`_, 

.. _`Docker Compose file for Orion-Broker`: https://github.com/smartfog/fogflow/tree/development/test/orion-ld
.. _`Docker Compose file for Scorpio Broker`: https://github.com/smartfog/fogflow/tree/development/test/scorpio
.. _`Docker Compose file for Stellio-context-broker`: https://github.com/smartfog/fogflow/tree/development/test/stellio-context-broker-ld
.. _`FogFlow installation Document`: https://fogflow.readthedocs.io/en/latest/setup.html


How to prepare and register an NGSI-LD based fog function for processing NGSI-LD data
================================================================================================




How to validate the entire workflow
================================================================================================





Using NGSI-LD specification implementation 
===============================================
Scorpio integration with FogFlow enable FogFlow task to communicate with scorpio Broker.
The figure below shows how data will transmit between scorpio broker, FogFlow broker and FogFlow task.

.. figure:: figures/scorpioIntegration.png

Integration steps
-----------------------

**Pre-Requisites:**

* FogFlow should be up and running with atleast one node.
* Scorpio Broker should be up and running.
* Create and trigger topology of two FogFunctions (`See Document`_).
* Create one fog Function (FogFunction-1) that publish update on FogFlow Broker (`Use template`_).
* Create another fog Function (FogFunction-2) that publish update on Scorpio Broker (`Use operator`_).

.. _`See Document`: https://fogflow.readthedocs.io/en/latest/intent_based_program.html.

.. _`Use template`: https://github.com/smartfog/fogflow/tree/development/application/template/NGSILD/python.

.. _`Use operator`: https://github.com/smartfog/fogflow/tree/development/application/operator/NGSI-LD-operator/NGSILDDemo.


**Below are the further steps for integration with Scorpio Broker.**

**Create any entity in Scorpio Broker**

.. code-block:: console

     curl -iX POST \
    'http://<Scorpio Broker>/ngsi-ld/v1/entities/' \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/ld+json' \
     -H 'Link: {{https://json-ld.org/contexts/person.jsonld}}; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
    -d '
        {
         "id": "urn:ngsi-ld:Vehicle:A13",
         "type": "Vehicle",
             "brandName": {
                  "type": "Property",
                  "value": "BMW",
                  "observedAt": "2017-07-29T12:00:04"
                },
                 "isParked": {
                   "type": "Relationship",
                   "object": "urn:ngsi-ld:OffStreetParking:Downtown",
                   "observedAt": "2017-07-29T12:00:04",
                    "providedBy": {
                        "type": "Relationship",
                        "object": "urn:ngsi-ld:Person:Bob"
                     	},
		}
        "location": {
                "type": "GeoProperty",
                "value": {
                        "type": "Point",
                        "coordinates": [-8.5, 41.2]
                }
        }
  }'



**FogFlow Will subscribe to scorpio Broker to get notification for every update to above created entity.**

.. code-block:: console

    curl -iX POST \
    'http://<Scorpio Broker>/ngsi-ld/v1/subscriptions/' \
      -H 'Content-Type: application/json' \
      -H 'Accept: application/ld+json' \
      -H 'Link: {{https://json-ld.org/contexts/person.jsonld}}; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
      -d '
      {
         "type": "Subscription",
         "entities": [{
                "id" : "urn:ngsi-ld:Vehicle:A13",
                "type": "Vehicle"
           }],
          "watchedAttributes": ["*"],
          "notification": {
                 "attributes": ["*"],
                  "format": "keyValues",
                 "endpoint": {
                        "uri": "http://<FogFLow Broker>/ngsi-ld/v1/notifyContext/",
                        "accept": "application/json"
                }
         }
    }'


**FogFlow Task will subscriber to FogFlow to get notification for furthur analysis.**

**NGSI-LD device will sends some update to scopio broker**

.. code-block:: console

    curl -iX PATCH \
    'http://<Scorpio Broker>/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A13/attrs' \
      -H 'Content-Type: application/json' \
      -H 'Accept: application/ld+json' \
      -H 'Link: {{https://json-ld.org/contexts/person.jsonld}}; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
      -d '
     {
	"brandName": {
		"type": "Property",
        	"value" : "BM2"
      		}
     }'



**Following process will occur internally in FogFLow**

* FogFunction-1 task will publish update on the FogFlow broker.
* FogFlow broker will send the notification to FogFunction-2 task.
* FogFunction-2 will convert this notification into scorpio update and send that update to scorpio broker.



Using NGSI-LD Adapter
===============================================


NGSI-LD Adapter is built to enable FogFlow Ecosystem to provide Linked Data to the users. `Scorpio Broker`_ being the first reference implementation of NGSI-LD Specification, is being used here for receiving the Linked-Data from Fogflow.

.. _`Scorpio Broker`: https://scorpio.readthedocs.io/en/latest/

The figure below shows how NGSI-LD Adapter works in transforming the NGSIv1 data from Fogflow into NGSI-LD data to Scorpio Broker.

.. figure:: figures/ngsi-ld-adapter.png

1. User sends a subscription request to the adapter. 
2. The adapter then forwards this request to the Fogflow broker, to subscribe itself for the Context Data specified in its request.
3. Context data update is received at Fogflow broker.
4. Adapter receives notification from the Fogflow broker for the subscribed data.
5. Adapter converts the received data into NGSI-LD data format and forwards it to the Scorpio broker. 


Running NGSI-LD Adapter
---------------------------

**Pre-Requisites:**

* Fogflow should be up and running with atleast one node.
* Scorpio broker should be up and running.

NGSI-LD Adapter can be run under Fogflow ecosystem using Fogflow Dashboard as given below. 

**Register an Operator:** Go to "Operator" in Operator Registry on Fogflow Dashboard. Register a new Operator with a Parameter Element as given below.
   
   Name: service_port ; Value: 8888
   
   (Is is assumed that the user has already gone through "REGISTER YOUR TASK OPERATORS" in `this`_ tutorial.)

.. _`this`: https://fogflow.readthedocs.io/en/latest/intent_based_program.html
   
**Register a Docker Image:** Go to "DockerImage" in Operator Registry and register an image fogflow/ngsildadapter:latest. Associate it with the above operator by choosing the operator from DropDown. Users can also build their image for NGSI-LD-Adapter by editing and running `build`_ file.

.. _`build`: https://github.com/smartfog/fogflow/blob/document-update/application/operator/NGSI-LD-Adapter/build

**Register a Fog Function** as shown in the figure below. In "SelectedType", provide the Entity Type (say "LD") of the Context Data that will be used to trigger this Fog Function. Choose the operator registered in Step#1 as the operator in Fog Function.

.. figure:: figures/fogfunction_ngsi-ld-adapter.png


**Trigger the Fog Function** by sending an update request to Fogflow Broker with the Entity Type as "LD" (or whatever is specified in Step#3 as the SelectedType). It should include fogflowIP and ngbIP in the attributes along with location metadata. Example request is given below:

.. code-block:: console

    curl -iX POST \
      'http://<Fogflow-Broker-IP>:8070/ngsi10/updateContext' \
      -H 'Content-Type: application/json' \
      -d '
      {
        "contextElements": [
        {
            "entityId": {
            "id": "LD001",
            "type": "LD",
            "isPattern": false
            },
            "attributes": [
                 {
                     "name": "fogflowIP",
                     "type": "string",
                     "value": "<IP>"
                 },
                 {
                     "name": "ngbIP",
                     "type": "string",
                     "value": "<IP>"
                 }
             ],
             "domainMetadata": [
                 {
                     "name": "location",
                     "type": "point",
                     "value": {
                                  "latitude": 52,
                                  "longitude": 67
                     }
                 }
             ]
        }
        ],
        "updateAction": "UPDATE"
       }'


NGSI-LD-Adapter task will be created and it will be listening on port 8888. Users can list it in the tasks running on either the cloud node or the edge node, whichever is nearest to the location provided in the metadata of the above request. 


How to use  NGSI-LD Adapter
-----------------------------

To use the NGSI-LD-Adapter for context data transformation, follow the below steps.


**Send subscription request** to LD-Adapter, it will forward the same request to Fogflow Broker. This is because the access to Fogflow broker will not be available directly to the user. Examle Subscription request is given below:

.. code-block:: console

    curl -iX POST \
      'http://<LD-Adapter-Host-IP>:8888/subscribeContext' \
      -H 'Content-Type: application/json' \
      -d '
    {
      "entities": [
        {
          "id": "Temperature.*",
          "type": "Temperature",
          "isPattern": true
        }
      ],
      "attributes": [
        "temp"
      ],
      "restriction": {
        "scopes": [
          {
            "scopeType": "circle",
            "scopeValue": {
              "centerLatitude": 49.406393,
              "centerLongitude": 8.684208,
              "radius": 2000
            }
          }
        ]
      },
      "reference": "http://<LD-Adapter-Host-IP>:8888"
    }'


**Send update request** to Fogflow Broker with an entity of type and attributes defined in the above subscription. An example request is given below:

.. code-block:: console

    curl -iX POST \
      'http://<Fogflow-Broker-IP>:8070/ngsi10/updateContext' \
      -H 'Content-Type: application/json' \
      -d '
      {
        "contextElements": [
          {
            "entityId": {
              "id": "Temperature001",
              "type": "Temperature",
              "isPattern": false
            },
            "attributes": [
              {
                "name": "temp",
                "type": "float",
                "value": 34
              }
            ],
            "domainMetadata": [
              {
              "name": "location",
              "type": "point",
              "value": {
                "latitude": 49.406393,
                "longitude": 8.684208
                }
              }
             ]
          }
        ],
        "updateAction": "UPDATE"
      }'


Check if the entity in NGSI-LD format has been updated on Scorpio Broker by visiting URL:  http://<Scorpio-Broker-IP:Port>/ngsi-ld/v1/entities?type=http://example.org/Temperature

Following code block shows the trasformed context data.

.. code-block:: console

    {"@context": ["https://schema.lab.fiware.org/ld/context", "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
    {"Temperature": "http://example.org/Temperature", "temp": "http://example.org/temp"}], "type": "Temperature", 
    "id": "urn:ngsi-ld:Temperature001", "temp": {"type": "Property", "value": 34}, "location": {"type": "GeoProperty", 
    "value": "{\"type\": \"point\", \"coordinates\": [49.406393, 8.684208]}"}}

package solarmodel

const AfterInvoiceIssuingEnvelope string = `[
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_header",
		"doc": "HEADER field of raw payload as emitted by the afterInvoiceIssuing solanteq RabbitMQ queue.",
		"fields": [
			{"name": "protocol", "type": {"type": "map", "values": "string"}, "doc": "Solar Protocol Information"},
			{"name": "messageId", "type": "string", "doc": "Solar message identifier."},
			{"name": "messageDate", "type": "string", "doc": "Date and time when the message was created."},
			{"name": "originator", "type": {"type": "map", "values": "string"}, "doc": "Information about the system of the source of the message."}		
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice_agreement",
		"doc": "Information about the contract to which the invoice document belongs inside SOLAR system.",
		"fields": [
			{"name": "id", "type": "int", "doc": "Identifier of the contract in the SOLAR system."},
			{"name": "number", "type": "string", "doc": "Number of contract (card_acceptor_id)."}	
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice_amount",
		"doc": "Information about the amount of the invoice.",
		"fields": [
			{"name": "amount", "type": "float", "doc": "Invoice amount"},
			{"name": "currency", "type": "string", "doc": "Numeric currency code according to ISO 4217."}		
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice_type",
		"doc": "Invoice type identification.",
		"fields": [
			{"name": "id", "type": "int", "doc": "Invoice type identifier in the SOLAR system."},
			{"name": "code", "type": "string", "doc": "Invoice type code."},
			{"name": "name", "type": "string", "doc": "Invoice type name."}		
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice_items_item_record_type",
		"doc": "Invoice Item Type identification.",
		"fields": [
			{"name": "id", "type": "int", "doc": "Invoice item type identifier in the SOLAR system."},
			{"name": "code", "type": "string", "doc": "Invoice item type code."},
			{"name": "name", "type": "string", "doc": "Invoice item type name."}		
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice_items_item_record_amount",
		"doc": "Information about the amount of the invoice item.",
		"fields": [
			{"name": "amount", "type": "float", "doc": "Invoice item mount"},
			{"name": "currency", "type": "string", "doc": "Numeric currency code according to ISO 4217."}		
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice_items",
		"doc": "Information about invoice items.",
		"fields": [
			{"name": "item", "type": "array", "doc": "List of items with information about one invoice item.",
			 "items": {					
					"name": "item_record",
					"doc": "Item information.",
					"type": "record",
					"fields": [
						{"name": "id", "type": "int", "doc": "Unique identifier of the invoice position in the SOLAR system."},
						{"name": "type", "type": "solar_afterInvoiceIssuing_body_invoice_items_item_record_type", "doc": "Invoice Item Type Identification."},	
						{"name": "itemAmount", "type": "solar_afterInvoiceIssuing_body_invoice_items_item_record_amount", "doc": "Information about the amount of the invoice item"},	
						{"name": "entriesCount", "type": "int", "doc": "Number of invoice records in position."}
					]
				}
			}	
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_invoice",
		"doc": "Invoice information.",
		"fields": [
			{"name": "id", "type": "int", "doc": "Invoice type identifier in the SOLAR system."},
			{"name": "type", "type": "solar_afterInvoiceIssuing_body_invoice_type", "doc": "Invoice type identification."},	
			{"name": "code", "type": "string", "doc": "Invoice code."},
			{"name": "agreement", "type": "solar_afterInvoiceIssuing_body_invoice_agreement", "doc": "Information about the contract to which the invoice document belongs inside SOLAR system."},
			{"name": "openingDate", "type": "string", "doc": "Invoice opening date."},
			{"name": "issueDate", "type": "string", "doc": "Invoice issue date."},
			{"name": "status", "type": "string", "doc": "Invoice status."},
			{"name": "invoiceAmount", "type": "solar_afterInvoiceIssuing_body_invoice_amount", "doc": "Information about the amount of the invoice."},
			{"name": "items", "type": "solar_afterInvoiceIssuing_body_invoice_items", "doc": "Information about invoice items."}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_job_type",
		"doc": "Process type identification.",
		"fields": [
			{"name": "id", "type": "int", "doc": "Process type identifier in the SOLAR system."},
			{"name": "code", "type": "string", "doc": "Process type code."},
			{"name": "name", "type": "string", "doc": "Process type name."}		
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body_job",
		"doc": "Solar information about the process by which the invoice is issued.",
		"fields": [
			{"name": "id", "type": "int", "doc": "Process identifier in the SOLAR system."},
			{"name": "type", "type": "solar_afterInvoiceIssuing_body_job_type", "doc": "Process type identification."}			
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing_body",
		"doc": "body field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "job", "type": "solar_afterInvoiceIssuing_body_job", "doc": "Solar information about the process by which the invoice is issued."},			
			{"name": "invoice", "type": "solar_afterInvoiceIssuing_body_invoice", "doc": "Invoice information."}			
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.payment_acquiring",
		"name": "solar_afterInvoiceIssuing",
		"doc": "Raw payload as emitted by the afterInvoiceIssuing solanteq RabbitMQ queue.",
		"fields": [
			{"name": "header", "type": "solar_afterInvoiceIssuing_header", "doc": "HEADER field of raw payload as emitted by the afterInvoiceIssuing solanteq RabbitMQ queue."},
			{"name": "body", "type": "solar_afterInvoiceIssuing_body", "doc": "BODY field of raw payload as emitted by the afterInvoiceIssuing solanteq RabbitMQ queue."}
		]
	}
]`

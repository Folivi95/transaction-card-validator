package solarmodel

const AfterProcessTxnEnvelope string = `[
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_header",
		"doc": "HEADER field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "protocol", "type": {"type": "map", "values": "string"}, "doc": "sample field doc"},
			{"name": "messageId", "type": "string", "doc": "sample field doc"},
			{"name": "messageDate", "type": "string", "doc": "sample field doc"},
			{"name": "originator", "type": {"type": "map", "values": "string"}, "doc": "sample field doc"},
			{"name": "receiver", "type": ["null", {"type": "map", "values": "string"}], "doc": "sample field doc"},
			{"name": "responseParams", "type": {
					"type": "record", "namespace": "co.saltpay.acquiring", "name": "solar_afterProcessTxn_header_responseParams",
					"fields": [
						{"name": "collectByReference", "type": "boolean", "doc": "sample field doc"},
						{"name": "includes", "type": {"type": "map", "values": {"type": "array", "items": {"type": "map", "values": ["string", "boolean"]}}}, "doc": "sample field doc"}
					]
				}, 
				"doc": "sample field doc"
			},
			{"name": "locale", "type": ["null", "string"], "doc": "sample field doc"},
			{"name": "workflow", "type": ["null", "string"], "doc": "sample field doc"},
			{"name": "preciseTime", "type": ["null", "string"], "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_transactionType",
		"doc": "body.txn.transactionType field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"},
			{"name": "directionClass", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_fees_fee_feeType",
		"doc": "body.txn.originator.fees.fee.feeType field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"},
			{"name": "feeClass", "type": "string", "doc": "sample field doc"},
			{"name": "directionClass", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_fees_fee_tariff_feeTariffValue",
		"doc": "body.txn.originator.fees.fee.tariff.feeTariffValue field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "amountType", "type": "string", "doc": "sample field doc"},
			{"name": "base", "type": {"type": "map", "values": ["string", "int"]}, "doc": "sample field doc"},
			{"name": "min", "type": {"type": "map", "values": ["string", "int"]}, "doc": "sample field doc"},
			{"name": "max", "type": {"type": "map", "values": ["string", "int"]}, "doc": "sample field doc"},
			{"name": "percentValue", "type": "float", "doc": "sample field doc"},
			{"name": "valueCurrency", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_fees_fee_tariff",
		"doc": "body.txn.originator.fees.fee.tariff field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "tariffType", "type": "string", "doc": "sample field doc"},
			{"name": "tariffClass", "type": "string", "doc": "sample field doc"},
			{"name": "tariffGroupId", "type": "long", "doc": "sample field doc"},
			{"name": "feeTariffValue", "type": "solar_afterProcessTxn_body_txn_originator_fees_fee_tariff_feeTariffValue", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_fees_fee",
		"doc": "body.txn.originator.fees.fee field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "feeType", "type": "solar_afterProcessTxn_body_txn_originator_fees_fee_feeType", "doc": "sample field doc"},
			{"name": "tariff", "type": "solar_afterProcessTxn_body_txn_originator_fees_fee_tariff", "doc": "sample field doc"},
			{"name": "direction", "type": "string", "doc": "sample field doc"},
			{"name": "amounts", "type": {"type": "map", "values": {"type": "map", "values": ["float", "string"]}}, "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_fees",
		"doc": "body.txn.originator.fees field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "fee", "type": {"type": "array", "items": {"type": "solar_afterProcessTxn_body_txn_originator_fees_fee"}}, "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_balances_balance_balanceType",
		"doc": "body.txn.originator.balances.balance.balanceType field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "long", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_balances_balance",
		"doc": "body.txn.originator.balances.balance field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "balanceType", "type": "solar_afterProcessTxn_body_txn_originator_balances_balance_balanceType", "doc": "sample field doc"},
			{"name": "currency", "type": "string", "doc": "sample field doc"},
			{"name": "own", "type": "int", "doc": "sample field doc"},
			{"name": "available", "type": "int", "doc": "sample field doc"},
			{"name": "blocked", "type": "int", "doc": "sample field doc"},
			{"name": "loan", "type": "int", "doc": "sample field doc"},
			{"name": "overlimit", "type": "int", "doc": "sample field doc"},
			{"name": "overdue", "type": "int", "doc": "sample field doc"},
			{"name": "creditLimit", "type": "int", "doc": "sample field doc"},
			{"name": "finBlocking", "type": "int", "doc": "sample field doc"},
			{"name": "interests", "type": "int", "doc": "sample field doc"},
			{"name": "penalty", "type": "int", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_balances",
		"doc": "body.txn.originator.balances field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "balance", "type": {"type": "array", "items": {"type": "solar_afterProcessTxn_body_txn_originator_balances_balance"}}, "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_system",
		"doc": "body.txn.originator.system field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"},
			{"name": "category", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_agreement",
		"doc": "body.txn.originator.agreement field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "number", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_accessor_type",
		"doc": "body.txn.originator.accessor.type field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator_accessor",
		"doc": "body.txn.originator.accessor field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "int", "doc": "sample field doc"},
			{"name": "number", "type": "string", "doc": "sample field doc"},
			{"name": "type", "type": "solar_afterProcessTxn_body_txn_originator_accessor_type", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_originator",
		"doc": "body.txn.originator field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "system", "type": "solar_afterProcessTxn_body_txn_originator_system", "doc": "sample field doc"},
			{"name": "agreement", "type": "solar_afterProcessTxn_body_txn_originator_agreement", "doc": "sample field doc"},
			{"name": "accessor", "type": "solar_afterProcessTxn_body_txn_originator_accessor", "doc": "sample field doc"},
			{"name": "number", "type": "string", "doc": "sample field doc"},
			{"name": "feeTotalAmount", "type": {"type": "map", "values": ["float", "string"]}, "doc": "sample field doc"},
			{"name": "fees", "type": "solar_afterProcessTxn_body_txn_originator_fees", "doc": "sample field doc"},
			{"name": "balances", "type": "solar_afterProcessTxn_body_txn_originator_balances", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_receiver_system_systemGroup",
		"doc": "body.txn.receiver.system.systemGroup field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "long", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_receiver_system",
		"doc": "body.txn.receiver.system field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "long", "doc": "sample field doc"},
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "name", "type": "string", "doc": "sample field doc"},
			{"name": "category", "type": "string", "doc": "sample field doc"},
			{"name": "systemGroup", "type": "solar_afterProcessTxn_body_txn_receiver_system_systemGroup", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_receiver_agreement",
		"doc": "body.txn.receiver.agreement field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "id", "type": "long", "doc": "sample field doc"},
			{"name": "number", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_receiver",
		"doc": "body.txn.receiver field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "system", "type": "solar_afterProcessTxn_body_txn_receiver_system", "doc": "sample field doc"},
			{"name": "agreement", "type": "solar_afterProcessTxn_body_txn_receiver_agreement", "doc": "sample field doc"},
			{"name": "accessor", "type": ["null", {"type": "map", "values": "string"}], "doc": "sample field doc"},
			{"name": "number", "type": "string", "doc": "sample field doc"},
			{"name": "feeTotalAmount", "type": ["null", {"type": "map", "values": ["double", "string"]}], "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_attributes_attribute",
		"doc": "body.txn.attributes.attribute field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "code", "type": "string", "doc": "sample field doc"},
			{"name": "attribute", "type": ["string", "long", {"type": "array", "items": ["null", "string"]}], "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn_attributes",
		"doc": "body.txn.attributes field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "attribute", "type": {"type": "array", "items": "solar_afterProcessTxn_body_txn_attributes_attribute"}, "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body_txn",
		"doc": "body.txn field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "txnId", "type": "long", "doc": "sample field doc"},
			{"name": "transactionType", "type": "solar_afterProcessTxn_body_txn_transactionType", "doc": "sample field doc"},
			{"name": "class", "type": "string", "doc": "sample field doc"},
			{"name": "category", "type": "string", "doc": "sample field doc"},
			{"name": "direction", "type": "string", "doc": "sample field doc"},
			{"name": "transactionDate", "type": "string", "doc": "sample field doc"},
			{"name": "settlementDate", "type": "string", "doc": "sample field doc"},
			{"name": "systemDate", "type": "string", "doc": "sample field doc"},
			{"name": "reference", "type": {"type": "map", "values": "string"}, "doc": "sample field doc"},
			{"name": "originator", "type": "solar_afterProcessTxn_body_txn_originator", "doc": "sample field doc"},
			{"name": "receiver", "type": "solar_afterProcessTxn_body_txn_receiver", "doc": "sample field doc"},
			{"name": "merchantData", "type": {"type": "map", "values": "string"}, "doc": "sample field doc"},
			{"name": "binTableData", "type": {"type": "map", "values": "string"}, "doc": "sample field doc"},
			{"name": "amounts", "type": {"type": "map", "values": {"type": "map", "values": ["double", "string"]}}, "doc": "sample field doc"},
			{"name": "attributes", "type": "solar_afterProcessTxn_body_txn_attributes", "doc": "sample field doc"},
			{"name": "txnConditions", "type": {"type": "map", "values": ["boolean", "string"]}, "doc": "sample field doc"},
			{"name": "txnDetails", "type": ["null", "string"], "doc": "sample field doc"},
			{"name": "response", "type": {"type": "map", "values": "string"}, "doc": "sample field doc"},
			{"name": "txnStatus", "type": "string", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn_body",
		"doc": "body field of raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "txn", "type": "solar_afterProcessTxn_body_txn", "doc": "sample field doc"}
		]
	},
	{
		"type": "record",
		"namespace": "co.saltpay.acquiring",
		"name": "solar_afterProcessTxn",
		"doc": "Raw payload as emitted by the afterProcessTxn solanteq RabbitMQ queue.",
		"fields": [
			{"name": "header", "type": "solar_afterProcessTxn_header", "doc": "sample field doc"},
			{"name": "body", "type": "solar_afterProcessTxn_body", "doc": "sample field doc"}
		]
	}
]`

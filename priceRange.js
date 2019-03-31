// iframe 中执行
var json = {"configurationData":[{"percentage":0,"shiptoCountry":"RU"},{"percentage":0,"shiptoCountry":"US"},{"percentage":0,"shiptoCountry":"CA"},{"percentage":0,"shiptoCountry":"ES"},{"percentage":0,"shiptoCountry":"FR"},{"percentage":0,"shiptoCountry":"UK"},{"percentage":0,"shiptoCountry":"NL"},{"percentage":0,"shiptoCountry":"IL"},{"percentage":0,"shiptoCountry":"BR"},{"percentage":0,"shiptoCountry":"CL"},{"percentage":0,"shiptoCountry":"AU"},{"percentage":0,"shiptoCountry":"UA"},{"percentage":0,"shiptoCountry":"BY"},{"percentage":0,"shiptoCountry":"JP"},{"percentage":0,"shiptoCountry":"TH"},{"percentage":0,"shiptoCountry":"SG"},{"percentage":0,"shiptoCountry":"KR"},{"percentage":0,"shiptoCountry":"ID"},{"percentage":0,"shiptoCountry":"MY"},{"percentage":0,"shiptoCountry":"PH"},{"percentage":0,"shiptoCountry":"VN"},{"percentage":0,"shiptoCountry":"IT"},{"percentage":0,"shiptoCountry":"DE"},{"percentage":0,"shiptoCountry":"SA"},{"percentage":0,"shiptoCountry":"AE"},{"percentage":0,"shiptoCountry":"PL"},{"percentage":0,"shiptoCountry":"TR"},{"percentage":0,"shiptoCountry":"PT"}],"configurationType":"percentage"};
$("#CHANNLE_FEE_1 > tbody > tr:nth-child(5) th").each((i,k) => {
	let percentage = parseFloat($(k).text());
	if (percentage < 0) {
		percentage = percentage / 1.5;
	}
	percentage = Math.round(percentage);
	json.configurationData[i].percentage = percentage;
});
json.configurationData = JSON.stringify(json.configurationData);
var s = JSON.stringify(json);
var template = `
var data = ${s};
$("#aeopNationalQuoteConfiguration").val(JSON.stringify(data));
clearData();
smtProAdd.initNationalData();	
`;
console.log(template);

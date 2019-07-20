// iframe 中执行
var json = {"configurationData":[{"absoluteQuoteMap":{"":0},"shiptoCountry":"RU"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"US"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"CA"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"ES"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"FR"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"UK"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"NL"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"IL"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"BR"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"CL"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"AU"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"UA"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"BY"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"JP"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"TH"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"SG"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"KR"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"ID"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"MY"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"PH"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"VN"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"IT"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"DE"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"SA"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"AE"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"PL"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"TR"},{"absoluteQuoteMap":{"":0},"shiptoCountry":"PT"}],"configurationType":"percentage"};
$("#CHANNLE_FEE_1 > tbody > tr:nth-child(5) th").each((i,k) => {
    let percentage = parseFloat($(k).text());
    if (percentage < 0) {
        percentage = percentage / 1.5;
    }
    percentage = Math.round(percentage);
    json.configurationData[i].absoluteQuoteMap = {"":percentage};
});
json.configurationData = JSON.stringify(json.configurationData);
var s = JSON.stringify(json);
var template = `
var data = ${s};
$("#aeopNationalQuoteConfiguration").val(JSON.stringify(data));
$('#setNationalProductPrice').find('.absoluteBox .modifyArea').html('');
smtProAdd.changePriceMode();
smtProAdd.initaeopNationalQuoteConfiguration('template');
`;
console.log(template);

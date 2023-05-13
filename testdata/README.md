```
curl -H "Content-Type: application/json" -H "X-Api-Key: $METEOCAT_API_KEY" https://api.meteo.cat/referencia/v1/municipis
```

```
curl -H "Content-Type: application/json" -H "X-Api-Key: $METEOCAT_API_KEY" https://api.meteo.cat/xema/v1/estacions/metadades?estat=ope&data=2023-03-11Z
```

```
cat testdata/metadades_totes_estacions.json| jq -r '.|sort_by(.nom)|.[]|([.nom, .codi])|@tsv'
```

```
curl -H "Content-Type: application/json" -H "X-Api-Key: $METEOCAT_API_KEY" https://api.meteo.cat/xema/v1/estacions/D5/metadades
```

```
curl -H "Content-Type: application/json" -H "X-Api-Key: $METEOCAT_API_KEY" https://api.meteo.cat/xema/v1/variables/mesurades/32/2023/03/12?codiEstacio=D5
```

https://api.meteo.cat/xema/v1/variables/mesurades/32/2017/03/27?codiEstacio=UG
# The New Quill: Tablas de procesos

## Nomenclatura

### Sentencia Lógica (SL)

Una **Sentencia Lógica** (SL) es la unidad mínima que el parser extrae del input del jugador. Tiene tres campos posicionales más un contenido opcional entrecomillado:

```
SL = (verbo, nombre, adjetivo, contenido-entrecomillado?)
```

Cualquiera de los tres primeros campos puede estar ausente — se representa con `_` (NullWord) en la cabecera del proceso que intente casar la SL.

A partir del `nombre`, el parser **infiere bindings**:

- Si el nombre resuelve a un Item declarado, `ITEM` queda ligado a ese Item.
- Si el nombre resuelve a un NPC declarado, `NPC` queda ligado a ese Character.

Estos bindings son visibles dentro del cuerpo del proceso que matchea la SL y, al anidar, también desde los sub-procesos.

### Cadena de SLs

El input del jugador se parte en **una o varias SLs encadenadas**, separadas por conectores (`y`, comas...). Cada SL se procesa secuencialmente por el ciclo de tablas.

Si una SL contiene un bloque entrecomillado, ese bloque queda reservado y **no se parsea** hasta que la tabla `npc` lo procesa (ver anidamiento más abajo).

### Tabla de procesos

Una tabla o lista que contiene una serie de procesos con expresiones de entrada. Dichas tablas son fijas y se ejecutan en un momento predeterminado.

Se termina la ejecución de las tablas con un `DONE`, o con un `OK` que es lo mismo que `DONE` pero muestra un mensaje de confirmación al usuario, el mensaje con label `ok`, o lo contrario `nook`:

| Action | Description                                                                                       |
|--------|---------------------------------------------------------------------------------------------------|
| `DONE` | Se termina silenciosamente la ejecución de las tablas, y se procesa la siguiente SL.              |
| `OK`   | Se imprime el mensaje `ok` y se termina la ejecución de las tablas, se procesa la siguiente SL.   |
| `NOOK` | Se imprime el mensaje `nook` y se termina la ejecución de las tablas (descartando resto de input).|
| `END`  | Termina la ejecución **solamente de la tabla actual** y se continúa con la siguiente tabla.       |

Caer al final del cuerpo de un proceso sin terminador explícito equivale a un `END` implícito: la tabla se da por terminada, se continua con la siguiente tabla si la hubiera.

### Proceso

Un proceso es una serie de instrucciones que se ejecutan cuando se cumple la condición de entrada en una tabla de procesos. Por ejemplo, en la tabla de `item`:

```plaintext
* *:
    NOT HERE ITEM
    SAY item_not_here, ITEM
    DONE

coger *:
    CARRIED ITEM
    SAY already_carried, ITEM
    NOOK

coger *:
    HERE vigilante
    SAY vigilante-here
    NOOK

// A partir de aquí sabemos que el objeto está y el vigilante no.

coger *:
    OVERWEIGHT ITEM
    SAY overweight, ITEM
    NOOK

coger *:
    GET ITEM
    SAY get, ITEM
    OK

examinar *:
    SAY ITEM.description
    DONE
```

Tabla de `npc`:

```plaintext
* *:
    NOT HERE NPC
    SAY npc_not_here, NPC
    DONE

buscar *:
    HERE NPC
    SAY "He encontrado a {1}.", NPC // Ver nota
    DONE

buscar *:
    SAY "No encuentro a {1}.", NPC
    DONE

examinar *:
    SAY NPC.description
    DONE
```

> Nota: Los mensajes no predefinidos en la sección de mensajes, crean un nuevo mensaje en la base de datos, con una
etiqueta única generada automáticamente. Se busca si el mensaje ya existe para evitar duplicados. Esto rompe la internacionalización (por definir), y se recomienda no hacerlo como regla general, pero aquí los usamos para mayor claridad.

## Tablas de procesos

Las tablas disponibles son:

| Orden | Tabla    | Ejecución                                                                |
|-------|----------|--------------------------------------------------------------------------|
| 0     | init     | Cuando arranca la aventura, o se reinicia                                |
| 1     | location | Cuando cambia de ubicación                                               |
| 2     | turn     | Tras el input, antes de `item`                                           |
| 3     | item     | Tras el input, sólo si la SL contiene un item                            |
| 4     | npc      | Tras el input, sólo si la SL contiene un NPC                             |
| 5     | response | Tras el input, en último lugar                                           |
| 6     | cron     | Por tiempo o turnos, procesos independientes                             |

No se pueden crear otras tablas. Se admiten `INCLUDE` para repartir el contenido entre varios ficheros.

Tras el input del jugador, se entra por la tabla 2 y se continúa con la 3, 4 y 5, salvo que se encuentre una instrucción `DONE`, `OK` o `NOOK`, en cuyo caso se detiene la ejecución de las tablas y se devuelve el control al jugador. Si se encuentra un `END`, se detiene la ejecución de la tabla actual y se continúa con la siguiente tabla compatible.

### Procesos

Un proceso es una entrada en una tabla de procesos que se ejecuta cuando se cumple la condición de entrada. Las condiciones de entrada son siempre **dos huecos**:

| Tabla    | Hueco 1                | Hueco 2                              | Ejemplos                                |
|----------|------------------------|--------------------------------------|-----------------------------------------|
| init     | `_`                    | `_`                                  | `_ _`                                   |
| location | localidad o `*`        | `_`                                  | `playa _`, `* _`                        |
| turn     | `EVERY` o `TIMEOUT`    | número entero positivo               | `EVERY 2`, `TIMEOUT 30`                 |
| item     | verbo, `*` o `_`       | noun de item, `*` o `_`              | `coger denario`, `coger *`, `* *`       |
| npc      | verbo, `*` o `_`       | noun de NPC, `*` o `_`               | `decir elfo`, `buscar *`                |
| response | verbo, `*` o `_`       | noun, `*` o `_`                      | `examinar playa`, `_ *`, `* *`          |
| cron     | unidad                 | número, o `"HH:MM:SS"` para `AT`     | `MINUTES 2`, `HOURS 1`, `AT "10:30:00"` |

#### Comodín y palabra vacía

- `*` — **comodín**, casa con cualquier label en esa posición.
- `_` — **NullWord**, indica que esa posición de la SL debe estar vacía para encajar.

Ambos se admiten en `init`, `location`, `item`, `npc` y `response`. **No** se admiten en `turn` ni `cron`, cuyas cabeceras codifican un disparador, no un patrón de matching.

#### Tabla de inicialización: `init`

Todas las cabeceras han de ser `_ _`, ya que todos los procesos se ejecutan en el arranque (o reinicio) y no hay SL que casar. Para agrupar procesos visualmente sólo se admiten comentarios.

#### Tabla de ubicaciones: `location`

El primer hueco siempre es una localidad o `*`, el segundo `_`. El proceso se ejecuta al entrar en esa localidad. Con `*`, el proceso aplica en todas.

#### Tabla de turnos: `turn`

Dos formas de cabecera:

- `EVERY n` — el proceso se ejecuta cada n turnos.
- `TIMEOUT n` — el proceso se ejecuta si el jugador no escribe en n segundos. El turno se da por consumido automáticamente.

Al consumirse un `TIMEOUT` se ejecuta este proceso y cualquier otro `EVERY` compatible con el nuevo número de turnos.

#### Tabla `item`

Se entra en esta tabla **sólo si la SL contiene un item** (el `nombre` resuelve a un Item declarado). Las cabeceras casan contra verbo + noun de item. El binding `ITEM` está vivo dentro de los procesos.

#### Tabla `npc`

Se entra en esta tabla **sólo si la SL contiene un NPC** (el `nombre` resuelve a un Character declarado). Las cabeceras casan contra verbo + noun de NPC. El binding `NPC` está vivo dentro de los procesos.

#### Tabla `response`

Tabla general, se ejecuta tras `turn`, `item` y `npc`. Admite `*` y `_` en cualquiera de los dos huecos.

```plaintext
encender linterna:
    SET linterna.on true
    SAY "Has encendido la linterna."
    OK

guardar partida:
    INPUT filename "Nombre de la partida: "
    NOT EMPTY filename
    SAVE filename
    OK

* insulto:
    SAY "{1} lo serás tú.", TITLECASE(SL.noun)
    QUIT

_ *:
    SAY "Cualquier nombre sin verbo."
    DONE

* _:
    SAY "Cualquier verbo sin nombre."
    DONE

* *:
    SAY "No te entiendo."
    DONE
```

#### Tabla `cron`

Los procesos se ejecutan por intervalos. El primer hueco es la unidad, el segundo el valor; con `AT` el segundo es una cadena con formato `"HH:MM:SS"`.

Unidades:

- `HOURS` — cada n horas (reloj de pared).
- `MINUTES` — cada n minutos.
- `SECONDS` — cada n segundos.
- `AT` — a una hora del día.

Los procesos `cron` son bloqueantes: si uno está en ejecución, el resto (cron u otra tabla) espera.

## Fallo de condición en un proceso

Cuando una condición no se cumple, el proceso se corta automáticamente y se intenta ejecutar el siguiente que encaje.

## Anidamiento de sub-procesos

Las tablas `response`, `item` y `npc` admiten **una capa** de anidamiento: dentro de un proceso pueden definirse sub-procesos con sus propias cabeceras. No se permite anidamiento en `init`, `location`, `turn` ni `cron`, ya que simplemente no tenemos _Sentencia Lógica_ disponible.

```plaintext
TABLE npc

    decir elfo:
        dar *:
            HAS NPC ITEM
            MOVE ITEM human
            SAY "Toma _.", ITEM
            OK

        dar *:
            SAY "No tengo _.", ITEM
            NOOK

        * *:
            SAY "El elfo no sabe responder a eso."
            NOOK

END TABLE
```

Reglas:

- **Una sola capa**: sin sub-sub-procesos.
- **Bindings heredados**: lo que armó el outer (`NPC`, `ITEM`, etc.) es visible desde los sub-procesos.
- **Preludio opcional**: el outer puede tener condactos antes de los sub-procesos. Esos condactos se ejecutan una vez al entrar al outer (y no se repiten en las re-entradas al mismo outer por sub-SLs distintas, ver más abajo).
- **Cierre implícito del outer**: termina al aparecer otra cabecera top-level o `END TABLE`.

## Hablar con NPCs: comillas y sub-SLs

Para dirigirse a un NPC el jugador **escribe la frase entre comillas**:

```
decir elfo "dame la llave"
```

El parser produce una SL principal con el contenido entrecomillado **sin parsear**:

```
SL = (verbo=decir, nombre=elfo, NPC=elfo, contenido="dame la llave")
```

Cuando el outer `decir elfo:` matchea en la tabla `npc`:

1. Se ejecuta el preludio del outer (si lo hay).
2. El parser extrae la **primera sub-SL** del contenido.
3. Los sub-procesos del outer se prueban contra esa sub-SL.
4. Si quedan más sub-SLs en el contenido, se **re-entra al mismo outer** con la siguiente sub-SL. **El preludio no se vuelve a ejecutar**; sólo los sub-procesos.

Ejemplo:

```
decir elfo "coge la llave y abre la puerta"
```

Expansión efectiva:

- SL principal: `decir elfo + contenido`.
- Sub-SL 1: `coger llave`.
- Sub-SL 2: `abrir puerta`.

Paso a paso:

1. Tabla `npc` matchea `decir elfo:`. Preludio ejecutado.
2. Sub-SL `coger llave` despachada contra los sub-procesos del outer.
3. Re-entrada al outer con sub-SL `abrir puerta`. Sólo sub-dispachar.

### SLs fuera de comillas

Lo que sigue al bloque entrecomillado son **SLs principales normales**, no dirigidas al NPC:

```
decir elfo "dame la llave" y soltar bolsa, quitar abrigo
```

- SL1: `decir elfo + contenido "dame la llave"` → entra a `npc`, dispatch interno.
- SL2: `soltar bolsa` → entra al ciclo `item/npc/response` con bindings limpios.
- SL3: `quitar abrigo` → entra al ciclo `item/npc/response` con bindings limpios.

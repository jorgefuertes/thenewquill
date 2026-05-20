# Condactos — Especificación

Documento vivo. Cada sección recoge las decisiones acordadas con el usuario.
Los puntos sin cerrar quedan marcados como `TBD`.

---

## 1. Vocabulario y alcance

**Linaje.** El motor se inspira en los sistemas clásicos de aventuras conversacionales
(PAWS de Spectrum, DAAD de PC, Aventuras AD), pero **no busca compatibilidad sintáctica
ni binaria** con ellos. Se conserva el principio rector — un mundo dirigido por
tablas de procesos, cada proceso compuesto por una secuencia de instrucciones
verbo+nombre+condición+acción — y se moderniza la sintaxis y la semántica donde
mejore la legibilidad o la potencia expresiva.

**Términos.**

- **Tabla de procesos**: contenedor con nombre que agrupa procesos. Cada tabla
  reservada tiene un disparador distinto en el ciclo de juego.
- **Proceso**: unidad ejecutable dentro de una tabla. Lleva una secuencia de
  condactos y, según el tipo de tabla, una cabecera (verbo+nombre, NPC, etc.).
- **Condacto**: pseudo-instrucción atómica que combina rol de condición y/o
  acción (de DAAD «condición + acción»). Cada condacto tiene un mnemónico y
  una firma de parámetros tipados.
- **SL** (sentence logic): el par verbo+nombre derivado del análisis de la
  entrada del jugador, usado por la tabla `response`.

**Alcance del MVP.**

- Cubrir las cinco tablas reservadas listadas en §3.
- Catálogo inicial de condactos suficiente para portar una aventura sencilla
  estilo "Aventura Original" (objetos, contenedores, localidades, NPCs,
  variables, mensajes).
- Compilación → DB serializable → runtime básico con interpretación de
  condactos.

**Fuera de alcance** (detalle en §11): debugger interactivo, traza de
ejecución, scripts externos, saltos arbitrarios, integración con motores
gráficos.

---

## 2. Sintaxis de la sección

> **Estado: propuesta**, pendiente de validación con el usuario.

### 2.1 Apertura de sección

Una sola directiva, alineada con el resto del lenguaje (`SECTION CONFIG`,
`SECTION ITEMS`, …):

```
SECTION TABLES
```

Se descarta `SECTION PROCESS TABLES` (manual antiguo): añadía una palabra sin
aportar contexto y rompía la regla léxica de "una palabra tras `SECTION`".

### 2.2 Delimitación de tabla

Cada tabla se abre con `TABLE <etiqueta>` y se cierra con `END TABLE`:

```
TABLE init
    ...
END TABLE
```

`<etiqueta>` es una de las cinco reservadas (§3). Una etiqueta no puede aparecer
dos veces en el mismo programa.

Se descarta la forma `BEGIN <KIND> TABLE` de la muestra antigua: redundante
(la etiqueta ya identifica el tipo) y verbosa.

### 2.3 Cuerpo de tabla

El cuerpo de toda tabla se divide en **procesos**. Cada proceso abre con una
**cabecera de dos tokens** terminada en `:` y su cuerpo son líneas de
condactos indentadas. El proceso queda implícitamente cerrado al encontrar
otra cabecera o `END TABLE` — no hay `END PROCESS`.

**Forma universal**: todas las tablas usan la misma estructura de cabecera —
dos tokens separados por un espacio, donde cada token puede ser una label
concreta, el comodín `*` o la NullWord `_`. Lo que esos dos tokens
*significan* depende del tipo de tabla:

| Tabla      | Slot 1                                    | Slot 2                                      |
|------------|-------------------------------------------|---------------------------------------------|
| `init`     | siempre `_`                               | siempre `_`                                 |
| `location` | label de localidad o `*`                  | siempre `_`                                 |
| `response` | label de verbo, `*` o `_`                 | label de nombre, `*` o `_`                  |
| `cron`     | unidad: `HOURS`, `MINUTES`, `SECONDS`, `TURNS` | número entero (cuenta de la unidad anterior) |
| `npc`      | label de nombre (de personaje) o `*`      | label de adjetivo, `*` o `_`                |

Casos especiales:

- **`init`** sólo admite `_ _` como cabecera. No se aceptan labels
  informativas; para anotar o agrupar visualmente usa comentarios.
- **`cron`** usa un patrón distinto a las demás tablas: slot 1 es una
  *unidad de tiempo* en mayúsculas, slot 2 es la *cuenta* en número entero.
  Ejemplos: `MINUTES 2:`, `TURNS 3:`, `SECONDS 30:`, `HOURS 1:`. La unidad
  `TURNS` se mide en turnos de juego; las otras tres se miden en tiempo
  real.
- **`location`** sólo discrimina por slot 1 (la localidad). Slot 2 queda
  reservado como `_` para mantener la forma uniforme. `*` en slot 1
  significa "este proceso corre en todas las localidades" (catch-all).
- **`npc`** identifica al personaje por el par (nombre, adjetivo) con que
  se declaró en `SECTION CHARACTERS` — la misma firma que usa el compilador
  para el dispatch.
- **`response`** casa contra la SL del jugador (slot 1 = verbo, slot 2 =
  nombre o adjetivo según la convención del lenguaje fuente).

Los comodines (§2.6) se aplican uniformemente en todas las tablas excepto
`cron` (donde ni `*` ni `_` tienen sentido — la cabecera codifica un
intervalo concreto).

### 2.4 La tabla `init`

`init` se ejecuta entera al arrancar. Su cuerpo es una lista de procesos
cuyas cabeceras son todas `_ _:` — `init` no discrimina, sólo agrupa
condactos en bloques sintácticos. **No se admiten labels informativas**;
para anotar usa comentarios:

```
TABLE init
    _ _:
        // World setup
        SET AT human start
        SET turns 0
        SET score 0
    _ _:
        // Light off, torch at the tomb
        DESTROY ant-on
        CREATE ant-off
        PUT_AT ant-off tumba
END TABLE
```

Todos los procesos `_ _:` corren en orden de aparición; el `END` o la
condición falsa de uno no afecta a los siguientes (sigue siendo la
política universal de §4.4). El programador puede tener un único proceso
gigante con todos los condactos o muchos pequeños — el resultado de la
ejecución es el mismo, pero la legibilidad cambia.

### 2.5 Líneas de condacto

Una línea es **un condacto** con sus parámetros, separados por espacios:

```
<mnemónico> <param1> <param2> ...
```

Las **condiciones** son condactos como cualquier otro: si una condición evalúa
falso, el resto del proceso se aborta y se sigue con el siguiente (§7).
Sintácticamente no hay distinción entre condiciones y acciones — todo es una
secuencia plana de condactos.

```
TABLE response
    take sword:
        CARRIED sword     # si ya lo tiene, condición falla, se aborta proceso
        HERE sword        # debe estar aquí
        GET sword
        OK
END TABLE
```

### 2.6 Comodines en las cabeceras

En las cabeceras de `init`, `location`, `response` y `npc` se admiten dos
labels especiales (ya reservadas en el sistema: existen como labels con IDs
fijos `1` y `2` en el `database` desde el arranque):

- `*` — **comodín**. Coincide con **cualquier** label en esa posición.
- `_` — **NullWord** (legado DAAD). Indica la **ausencia** de palabra en esa
  posición — útil para `response` (verbo o nombre ausentes en la SL del
  jugador) y obligatorio en los slots reservados de `init` y slot 2 de
  `location`.

Ejemplos:

```
TABLE response
    coger *:        // "coger <cualquier cosa>"
        ...
    * denario:      // "<cualquier verbo> denario"
        ...
    * *:            // catch-all: lo que no encajó arriba
        WRITE i-dont-understand
    examinar _:     // "examinar" sin nombre
        ...

TABLE location
    * _:            // se ejecuta en cualquier localidad
        WRITE turn-counter

TABLE npc
    * _:            // cualquier personaje sin adjetivo
        WRITE generic-greeting
    * *:            // cualquier personaje, cualquier adjetivo (más amplio)
        WRITE generic-snub

TABLE init
    _ _:            // único formato admitido en init
        SET turns 0
```

**Orden de resolución.** Cuando hay varios procesos cuya cabecera podría
coincidir (por wildcards), el runtime los prueba **en orden de aparición**
en el fuente y aborta el siguiente sólo si el actual ejecuta un condacto
que cierra el turno (`DONE`/`OK`/etc., concretado en §6.3). Esto da control
explícito al programador: "lo específico primero, lo genérico al final".

`*` y `_` **no** se admiten en cabeceras de `cron`: ahí la cabecera codifica
un intervalo concreto (unidad + cuenta), no un patrón de matching.

### 2.7 Comentarios, includes y blancos

Reutilizamos las convenciones léxicas ya existentes en el lenguaje:

- Comentarios de línea con `//` y de bloque con `/* ... */`.
- Líneas en blanco se ignoran.
- `INCLUDE "fichero.inc.adv"` se admite dentro de la sección, dividiendo una
  tabla larga en varios ficheros si conviene. Las `TABLE ... END TABLE`
  pueden cruzar fronteras de fichero, pero no es buena práctica: lo habitual
  es **una tabla por fichero include**.

### 2.8 Ejemplo completo

```
SECTION TABLES

TABLE init
    _ _:
        // World setup
        MOVE human start
    _ _:
        // Counters
        SET score 0
        SET hunger 0
END TABLE

TABLE location
    cave _:
        ZERO visited-cave            // sólo la primera vez
        SET visited-cave 1
        WRITELN seen-cave-first
    forest _:
        WRITELN seen-forest
    * _:
        WRITELN you-feel-the-wind    // en cualquier localidad
END TABLE

TABLE response
    coger denario:
        HERE denario
        NOT CARRIED denario
        GET denario
        OK
    dejar denario:
        CARRIED denario
        DROP denario
        OK
    coger *:                         // "coger" cualquier otra cosa
        WRITELN cannot-take
    examinar _:                      // "examinar" sin objeto
        LOOK
        DONE
    * *:                             // catch-all
        WRITELN i-dont-understand
END TABLE

TABLE cron
    SECONDS 30:
        CHANCE 25
        WRITELN thunder-rumble
    MINUTES 5:
        INC hunger 1
        GT hunger 10
        WRITELN hungry-warning
    TURNS 1:
        INC turns 1
END TABLE

TABLE npc
    sage _:                          // personaje "sage" declarado como `sage _`
        WRITELN sage-greeting
    troll _:
        ATTACK
    * _:                             // cualquier otro NPC (adj `_`)
        WRITELN generic-npc-snub
    * *:                             // fallback aún más amplio
        WRITELN generic-npc-snub
END TABLE
```

---

## 3. Tablas reservadas

Hay un conjunto cerrado de tablas reservadas. Cada una se identifica por su
**etiqueta canónica** y dispara en un momento concreto del ciclo de juego.
No hay tablas definidas por el usuario fuera de este conjunto en el MVP.

| Etiqueta   | Disparador                                                                       |
|------------|----------------------------------------------------------------------------------|
| `init`     | Se ejecuta una sola vez al arrancar la aventura, antes del primer turno.         |
| `location` | Se ejecuta al entrar en una nueva localidad o al pedir su descripción.           |
| `response` | Se ejecuta al procesar la entrada del jugador, casando con la SL (verbo+nombre). |
| `cron`     | Se ejecuta en segundo plano según un programa temporal (ver más abajo).          |
| `npc`      | Se ejecuta cuando el jugador interactúa con un NPC (hablar/atacar/etc.).         |

**Notas.**

- Se descarta la tabla `turn` que aparecía en el manual: su rol queda cubierto
  por `cron` (para acciones periódicas) y por `response` (para reacciones a
  acciones del jugador). Si más adelante se demuestra necesaria, se reconsidera.
- `cron` programa cada entrada con una cabecera de dos tokens
  `<unidad> <cuenta>:` donde `<unidad>` es una de las cuatro palabras
  reservadas y `<cuenta>` es un número entero positivo:

  | Cabecera        | Significado            |
  |-----------------|------------------------|
  | `SECONDS 15:`   | cada 15 segundos       |
  | `MINUTES 5:`    | cada 5 minutos         |
  | `HOURS 1:`      | cada hora              |
  | `TURNS 10:`     | cada 10 turnos         |

  `HOURS`, `MINUTES` y `SECONDS` son reloj de pared; `TURNS` se mide en
  turnos completos del jugador (§7.1). No se admiten otras unidades ni
  combinaciones (`MINUTES 4 SECONDS 30` no existe; se aproxima con
  `SECONDS 270` o con dos procesos `cron` distintos).
- `npc`: cada proceso lleva en su cabecera el par **nombre-adjetivo** del
  personaje al que aplica (igual que un personaje se declara en
  `SECTION CHARACTERS` con un noun y un adjetivo). Esto permite indexar
  por personaje y dispatch directo en runtime.
- Exclusión mutua: si una tabla está ejecutándose, las demás esperan. La
  regla detallada (qué bloquea a qué, qué pasa si se acumulan disparos
  de `cron`) se fija en §7.

---

## 4. Anatomía del proceso

> **Estado: propuesta**, pendiente de validación.

Un **proceso** es la unidad ejecutable dentro de una tabla. Tiene tres partes:
cabecera (cuando aplica), cuerpo, y un cierre implícito.

### 4.1 Cabecera

La forma léxica de la cabecera se describe en §2.3 y §2.6 (dos tokens, cada
uno label o comodín). Conceptualmente, la cabecera es la **clave de
dispatch** que decide cuándo este proceso es elegible para ejecutarse:

| Tabla      | La cabecera identifica…                                                  |
|------------|--------------------------------------------------------------------------|
| `init`     | nada (`_ _` siempre); todos los procesos son elegibles al arrancar       |
| `location` | la localidad sobre la que aplica el proceso (slot 1; slot 2 `_`)         |
| `response` | el par verbo+nombre de la SL del jugador                                 |
| `cron`     | el intervalo: unidad (`HOURS`/`MINUTES`/`SECONDS`/`TURNS`) + cuenta      |
| `npc`      | el par nombre+adjetivo del personaje (la misma firma de declaración)     |

La resolución contra wildcards (`*`, `_`) se hace en orden de aparición en el
fuente (§2.6). Esto permite escribir "lo específico primero, lo genérico al
final".

### 4.2 Cuerpo

El cuerpo es una **secuencia ordenada de condactos**. Cada línea es un
condacto con sus parámetros (§2.5). No hay distinción sintáctica entre
condiciones y acciones: todo es un condacto y se evalúa en orden.

### 4.3 Cierre implícito

El proceso queda cerrado cuando aparece otra cabecera al mismo nivel o
`END TABLE`. **No hay `END PROCESS`**. La indentación es indicativa pero no
significativa: el cierre se detecta por la siguiente cabecera, no por
desindentación.

### 4.4 Resultado de evaluar un condacto

Cada condacto, al ejecutarse, produce uno de tres resultados:

- **Sigue**: la ejecución continúa con el siguiente condacto del mismo proceso.
- **Aborta proceso** (≈ "condición falsa"): se detiene **este** proceso y el
  motor pasa al siguiente proceso elegible de la **misma tabla**.
- **Aborta turno**: se detiene este proceso **y la tabla entera**, y no se
  procesarán más tablas en este turno. Excepción: las tablas que ya empezaron
  a ejecutarse en este disparo terminan su proceso actual; las que aún no
  habían empezado quedan canceladas hasta el próximo turno.

Las condiciones "puras" (`AT`, `HERE`, `CARRIED`, `EQ VAR`, etc.) son
condactos que devuelven *aborta proceso* cuando la comprobación es falsa, y
*sigue* cuando es verdadera. Las acciones "puras" siempre devuelven *sigue*.
Los terminadores (`END`, `DONE`, `OK`) están especificados en §6.3.

### 4.5 Negación

El prefijo léxico **`NOT`** invierte el resultado de la siguiente condición:

```
NOT CARRIED sword     // pasa si el jugador NO lleva la espada
NOT AT cave           // pasa si NO está en la cueva
```

`NOT` no es un mnemónico — es un **modificador léxico** independiente del
catálogo de condactos. La forma general de una línea es:

```
[NOT] <mnemónico> [<param1> <param2> ...]
```

Es decir: si el primer token de la línea es `NOT`, el siguiente es el
mnemónico y el primero arma el bit `NOT` del `flags` byte del condacto en el
bytecode (§8.3). Si el primer token no es `NOT`, entonces él mismo es el
mnemónico. Esto encaja con la regla de "mnemónicos en una sola palabra"
(§6.1) sin ambigüedad.

`NOT` es **palabra reservada** del lenguaje y no puede usarse como label
(de localidad, item, mensaje, etc.).

`NOT` se aplica sólo a la condición inmediatamente siguiente, no acumula
(`NOT NOT CARRIED sword` es error de compilación, no doble-negación).
Aplicarlo a una acción es también un error de compilación (no tiene
sentido negar un efecto).

### 4.6 Restricciones por tabla

Algunos condactos sólo tienen sentido en ciertas tablas. Las restricciones
duras (validadas en compilación) se enumeran en §9. Por ahora:

- `EVERY <intervalo>:` sólo aparece como **cabecera** de procesos en `cron`,
  nunca como condacto.
- Las cabeceras informativas de `init` (§2.4) no son condactos.

---

## 5. Parámetros

> **Estado: propuesta**, pendiente de validación.

Un condacto declara una **firma de parámetros tipados**. Los tipos disponibles
parten del enum `Param` ya definido en `internal/adventure/process/param.go`,
con un par de añadidos propuestos:

### 5.1 Tipos

| Tipo       | Significado                                         | Reconocimiento léxico                            |
|------------|-----------------------------------------------------|--------------------------------------------------|
| `Str`      | cadena literal                                       | `"texto entre comillas"`                         |
| `Num`      | entero (int32)                                       | `42`, `-3`, `0`                                  |
| `Bool`     | booleano                                             | `true`, `false`                                  |
| `WordID`   | referencia a palabra del vocabulario                 | label nuda: `north`                              |
| `LocID`    | referencia a localidad                               | label nuda: `cave`                               |
| `VarID`    | referencia a variable                                | label nuda: `score`                              |
| `ItemID`   | referencia a objeto                                  | label nuda: `sword`                              |
| `MsgID`    | referencia a mensaje                                 | label nuda: `seen-cave`                          |
| `TableID`  | referencia a tabla                                   | label nuda: `response`                           |
| `ProcID`   | referencia a proceso dentro de una tabla             | label nuda; resolución en §8                     |
| `CharID`   | **(nuevo)** referencia a personaje (NPC o humano)   | label nuda: `sage`                               |
| `BlobID`   | **(nuevo)** referencia a recurso binario (gfx/snd)  | label nuda: `logo`                               |

Los dos tipos nuevos (`CharID`, `BlobID`) son necesarios para condactos como
`IS_HUMAN char`, `IN_ROOM char loc` o `SHOW <blob>`. Pendiente: añadir las
constantes al enum `Param` cuando llegue la implementación.

### 5.2 Reconocimiento por contexto

Una **label nuda** (sin comillas, sin números, sin sufijos) es ambigua a nivel
léxico — puede ser cualquier tipo `*ID`. La firma del condacto resuelve la
ambigüedad: la posición del parámetro determina qué tipo se espera, y el
compilador comprueba que la label apunta a un registro del kind correcto.

Ejemplo. Para el condacto `AT <LocID>`:

```
AT cave        // ok: 'cave' existe y es una localidad
AT sword       // error de compilación: 'sword' es Item, se esperaba LocID
AT unknown     // error de compilación: label no definida
```

### 5.3 Variantes y conversiones

- **`Num` vs `Str`**: si una firma espera `Num` y el token es un literal entero
  o decimal, se acepta. Una cadena entre comillas en posición `Num` es error.
- **`Bool`**: se admite estrictamente `true` o `false`. No se infiere desde
  `Num`.
- **`*ID` y comodines**: `*` y `_` (§2.6) sólo son válidos en **cabeceras de
  proceso**, no como parámetros de condacto. Un condacto que necesite "ningún
  objeto" debe usar otra forma (por ejemplo, ausencia del condacto).
- **Listas / cantidades variables**: el MVP **no soporta varargs en general**.
  La única excepción son `WRITE` y `WRITELN` (§6.3), donde el número de
  argumentos viene determinado por la cantidad de placeholders `_` del
  mensaje referenciado. Esa aridad se conoce en compilación (basta con leer
  el mensaje correspondiente) y queda fija para cada call site — no es
  realmente "variable", sólo "dependiente del mensaje".

### 5.4 Validación en compilación

Cada parámetro pasa por dos comprobaciones:

1. **Léxica**: el token tiene la forma correcta para alguno de los tipos
   primitivos (`"..."` para `Str`, dígitos para `Num`, etc.) **o** es una
   label válida según el regex existente del proyecto.
2. **Semántica**: si la firma espera `*ID`, se resuelve la label contra los
   stores correspondientes (`Locations.Get().WithLabel(...)` para `LocID`,
   etc.) y se asegura que el kind del registro coincide.

Errores típicos detectables en compilación:

- Aridad incorrecta (faltan o sobran parámetros).
- Tipo primitivo incompatible (`SET VAR score "hola"` cuando se espera `Num`).
- Label no resuelta (ningún registro con esa label).
- Kind incorrecto (la label existe pero es de otro tipo).

Los errores de runtime se documentan en §9.

### 5.5 Serialización

Los detalles binarios del bytecode se especifican en §8. Conceptualmente, cada
parámetro se serializa como **byte de tipo** (corresponde al enum `Param`) +
**payload** (uint32 para `*ID`, longitud + bytes para `Str`, etc.).

---

## 6. Catálogo

### 6.1 Condiciones

> **Estado: propuesta**, pendiente de validación.

Una **condición** es un condacto que sólo decide si el proceso continúa: si
la condición es cierta, sigue; si es falsa, aborta el proceso (§4.4). Las
condiciones se pueden prefijar con `NOT` (§4.5) para invertir su resultado.

**Convenciones para esta pasada (confirmadas):**

- Mnemónicos siempre en **una sola palabra**, MAYÚSCULAS, con guión bajo
  donde mejora la legibilidad (`IN_ROOM`). `NOT` no es un mnemónico, es
  modificador léxico (§4.5).
- El **humano** se referencia con la label fija `human` — el personaje
  marcado con `is human` en `SECTION CHARACTERS` recibe ese alias en
  compilación. Las condiciones aceptan cualquier `CharID`.
- "Es NPC" se expresa como `NOT HUMAN <char>`; no se añade `IS_NPC`.
- Para un `Item` con `Created == false` (no ha entrado en juego), las
  condiciones de localización (`HERE`, `CARRIED`, etc.) devuelven **falso**:
  un objeto no creado no está en ningún sitio.

**Rangos de opcodes:**

- `0x01..0x3F` — condiciones
- `0x40..0xBF` — acciones de mundo
- `0xC0..0xFF` — salida y control de flujo

#### Localización e inventario

Las condiciones de visibilidad/posesión son **recursivas a través de
contenedores**. El modelo añade un campo `Open bool` al `Item` (§8.2): un
contenedor abierto deja ver su contenido, uno cerrado no.

| Opcode | Mnemónico | Firma             | Semántica                                          |
|--------|-----------|-------------------|----------------------------------------------------|
| `0x01` | `AT`      | `<LocID>`         | El humano está en esa localidad.                   |
| `0x02` | `HERE`    | `<ItemID>`        | El item está visible en la localidad actual: aquí directamente, o dentro de una cadena de contenedores abiertos que están aquí. |
| `0x03` | `CARRIED` | `<ItemID>`        | El humano lleva el item: directamente o dentro de cualquier contenedor que él lleve (la transparencia del contenedor **no afecta** — si lo llevas, lo llevas). |
| `0x04` | `WORN`    | `<ItemID>`        | El humano lleva el item puesto. `item.Worn == true`. No recursivo. |
| `0x05` | `PRESENT` | `<ItemID>`        | Equivalente a `HERE OR CARRIED`, con las semánticas de cada uno. |
| `0x06` | `EXIST`   | `<ItemID>`        | El item ha sido creado (está "en juego"). `item.Created == true`. |
| `0x07` | `IN_ROOM` | `<CharID> <LocID>`| El personaje (NPC o humano) está en esa localidad. |
| `0x0F` | `IS_OPEN` | `<ItemID>`        | El item es contenedor y está abierto. `item.Container && item.Open`. |

**Asimetría intencional `HERE` vs `CARRIED`.** `HERE` modela *visibilidad* —
lo que el jugador podría ver en la sala — y por tanto respeta el flag `Open`
de los contenedores intermedios. `CARRIED` modela *posesión física*: si llevas
una mochila cerrada con una moneda dentro, sigues "llevando" la moneda aunque
no la veas. `PRESENT` es la disyunción literal de ambos y hereda esa
asimetría: una moneda en una mochila cerrada que llevas es `PRESENT`
(porque está `CARRIED`), aunque no esté `HERE`.

Si alguna aventura prefiere otra política (e.g. "PRESENT = visible siempre"),
se ajusta en una versión futura. El MVP fija la asimetría arriba.

**Recursividad y ciclos.** El recorrido de la cadena de contenedores se hace
hasta encontrar localidad, personaje o item no-contenedor. Se detecta y se
rechaza un ciclo (error de runtime, §9.2) — no debería ocurrir si los
condactos de movimiento (§6.2) validan correctamente.

Ejemplo:

```
TABLE response
    take sword:
        HERE sword                  // visible aquí, incl. contenedores abiertos
        NOT CARRIED sword           // y no lo llevo ya
        GET sword
        OK
    open chest:
        HERE chest
        NOT IS_OPEN chest           // no abrirlo dos veces
        OPEN chest
        OK
```

#### Variables

| Opcode | Mnemónico | Firma             | Semántica                              |
|--------|-----------|-------------------|----------------------------------------|
| `0x08` | `EQ`      | `<VarID> <Num>`   | `var == num`.                          |
| `0x09` | `NE`      | `<VarID> <Num>`   | `var != num`.                          |
| `0x0A` | `GT`      | `<VarID> <Num>`   | `var > num`.                           |
| `0x0B` | `LT`      | `<VarID> <Num>`   | `var < num`.                           |
| `0x0C` | `ZERO`    | `<VarID>`         | `var == 0`. Atajo de `EQ <var> 0`.     |

Ejemplo:

```
TABLE cron
    EVERY 1m:
        INC hunger 1
        GT hunger 10
        WRITE hungry-warning
```

#### Personajes

| Opcode | Mnemónico | Firma          | Semántica                                       |
|--------|-----------|----------------|-------------------------------------------------|
| `0x0D` | `HUMAN`   | `<CharID>`     | El personaje es el jugador (`Human == true`).   |

Para "es un NPC" se usa `NOT HUMAN char`. No se duplica el catálogo.

Ejemplo:

```
TABLE npc
    sage:
        NOT HUMAN sage         // siempre cierto en esta tabla, pero
                               // documenta intención y blinda ante refactors
        WRITE sage-greeting
```

#### Aleatorio

| Opcode | Mnemónico | Firma     | Semántica                                              |
|--------|-----------|-----------|--------------------------------------------------------|
| `0x0E` | `CHANCE`  | `<Num>`   | Pasa con probabilidad `num`% (0–100, fuera de rango = error de compilación). |

Ejemplo:

```
TABLE cron
    EVERY 30s:
        CHANCE 25
        WRITE thunder-rumble
```

#### Resumen §6.1

15 condiciones cubren el MVP:
`AT`, `HERE`, `CARRIED`, `WORN`, `PRESENT`, `EXIST`, `IN_ROOM`, `IS_OPEN`,
`EQ`, `NE`, `GT`, `LT`, `ZERO`, `HUMAN`, `CHANCE`.

Opcodes asignados: `0x01..0x0F`. Quedan libres `0x10..0x3F` (48 huecos) para
condiciones futuras (e.g. `CONTAINS`, `WORN_BY`, `ADJACENT`, `BIGGER_THAN`).

**Dependencias generadas por §6.1 sobre otras secciones:**

- §8.2 — el struct `Item` gana un campo `Open bool` para soportar
  `IS_OPEN` y la recursividad de `HERE`/`PRESENT`.
- §6.2 — el catálogo de acciones añadirá `OPEN <ItemID>` y `CLOSE <ItemID>`
  como mutadores del flag.

### 6.2 Acciones de mundo

> **Estado: propuesta**, pendiente de validación.

Una **acción de mundo** muta el estado del juego (posiciones, flags, valores
de variables) y siempre devuelve "sigue" — no aborta el proceso. Aplicarles
`NOT` es error de compilación (§4.5).

**Convenciones:**

- Las acciones son **raw**: no llevan precondiciones implícitas. Si necesitas
  validar (e.g. "el item está aquí antes de cogerlo"), encadénalo con
  condiciones explícitas. Esto es predecible y facilita escribir mensajes de
  error específicos.
- El humano se referencia con la label `human` (§6.1).
- "Localidad actual" del humano se consulta del estado del intérprete (§7.1)
  — no se pasa como parámetro.

#### Movimiento de entidades

| Opcode | Mnemónico  | Firma                | Semántica                                                       |
|--------|------------|----------------------|-----------------------------------------------------------------|
| `0x40` | `MOVE`     | `<CharID> <LocID>`   | Coloca el personaje en la localidad. Sobreescribe su `LocationID`. **Sólo personajes**: items usan `PUT_AT` / `GIVE` / `PUT_INTO`. |
| `0x41` | `PUT_AT`   | `<ItemID> <LocID>`   | Coloca el item en la localidad. `item.At = loc.ID`.             |
| `0x42` | `GIVE`     | `<ItemID> <CharID>`  | Pone el item en el inventario directo del personaje. `item.At = char.ID`. |
| `0x43` | `PUT_INTO` | `<ItemID> <ItemID>`  | Mete el primer item dentro del segundo (que debe ser contenedor). `item.At = container.ID`. |
| `0x44` | `GOTO`     | `<LocID>`            | Atajo: equivalente a `MOVE human <loc>`. Mueve al jugador.       |

Notas:

- `PUT_INTO` valida en compilación que el segundo argumento corresponde a un
  item con `Container == true` (§9.1).
- Tras `MOVE human <loc>` o `GOTO <loc>`, si la localidad cambió, el
  intérprete dispara la tabla `location` al terminar la respuesta actual
  (§7.2).
- Si más adelante hace falta intercambiar dos items sin pisar destinos
  (estilo `SWAP` de DAAD), se añadirá como condacto futuro. No forma parte
  del MVP.

#### Inventario del jugador (atajos)

| Opcode | Mnemónico  | Firma         | Semántica                                                   |
|--------|------------|---------------|-------------------------------------------------------------|
| `0x50` | `GET`      | `<ItemID>`    | Atajo: `GIVE <item> human`. El jugador toma el item.        |
| `0x51` | `DROP`     | `<ItemID>`    | Atajo: `PUT_AT <item> <localidad actual>`. El jugador suelta. |
| `0x52` | `WEAR`     | `<ItemID>`    | Marca el item como puesto: `item.Worn = true`. No cambia `At`. |
| `0x53` | `UNWEAR`   | `<ItemID>`    | Marca el item como no-puesto: `item.Worn = false`.          |

Notas:

- `WEAR` no exige que el item esté `CARRIED` previamente — eso es
  responsabilidad del programador encadenarlo. La validación en compilación
  comprueba que el item tiene `Wearable == true`.
- `UNWEAR` es la operación opuesta a `WEAR`. Se usa este nombre (en lugar
  del clásico DAAD `REMOVE`) para evitar la ambigüedad con "remove from
  world" — para eso está `DESTROY`.

#### Contenedores

| Opcode | Mnemónico | Firma         | Semántica                                          |
|--------|-----------|---------------|----------------------------------------------------|
| `0x60` | `OPEN`    | `<ItemID>`    | Marca el contenedor como abierto. `item.Open = true`. |
| `0x61` | `CLOSE`   | `<ItemID>`    | Marca el contenedor como cerrado. `item.Open = false`. |

Notas:

- Validación en compilación: el item debe tener `Container == true`.

#### Existencia (entrada y salida del juego)

| Opcode | Mnemónico | Firma         | Semántica                                          |
|--------|-----------|---------------|----------------------------------------------------|
| `0x70` | `CREATE`  | `<ItemID>`    | El item entra en juego. `item.Created = true`. Su posición se mantiene (la que tuviera por defecto o por un PUT_AT previo). |
| `0x71` | `DESTROY` | `<ItemID>`    | El item sale de juego. `item.Created = false`. Las condiciones de localización lo verán como ausente. |

`DESTROY` no borra el item de la base de datos — sólo apaga el flag. Esto
permite re-`CREATE`-arlo más tarde si el guión lo requiere.

#### Mutación de variables

| Opcode | Mnemónico | Firma             | Semántica                              |
|--------|-----------|-------------------|----------------------------------------|
| `0x80` | `SET`     | `<VarID> <Num>`   | Asigna `var = num` (entero con signo). |
| `0x81` | `INC`     | `<VarID> <Num>`   | `var += num` (num ≥ 0).                |
| `0x82` | `DEC`     | `<VarID> <Num>`   | `var -= num` (num ≥ 0).                |

Notas:

- `Num` es `int32`. Desbordamiento en `INC`/`DEC` es error de runtime (§9.2)
  — no se silencian wraps.
- `INC` y `DEC` exigen `num ≥ 0` en compilación. Para incrementos negativos
  se usa `DEC` con el valor absoluto. Esto evita la confusión de signos.

#### Resumen §6.2

16 acciones de mundo:
`MOVE`, `PUT_AT`, `GIVE`, `PUT_INTO`, `GOTO`,
`GET`, `DROP`, `WEAR`, `UNWEAR`,
`OPEN`, `CLOSE`,
`CREATE`, `DESTROY`,
`SET`, `INC`, `DEC`.

Opcodes asignados: `0x40..0x44`, `0x50..0x53`, `0x60..0x61`, `0x70..0x71`,
`0x80..0x82`. El rango `0x40..0xBF` permite hasta 128 acciones — quedan
~112 huecos para el futuro (e.g. `SWAP`, `LET <var> <expr>`,
operaciones de string, etc.).

**Ejemplo completo:**

```
TABLE response
    take sword:
        HERE sword
        NOT CARRIED sword
        GET sword
        OK
    open chest:
        HERE chest
        NOT IS_OPEN chest
        OPEN chest
        OK
    put key into chest:
        CARRIED key
        IS_OPEN chest
        PUT_INTO key chest
        OK
    wear cloak:
        CARRIED cloak
        NOT WORN cloak
        WEAR cloak
        OK
    take off cloak:
        WORN cloak
        UNWEAR cloak
        OK
END TABLE
```

### 6.3 Salida y control

> **Estado: propuesta**, pendiente de validación.

Los condactos de **salida** producen efectos visibles (texto, imagen, audio)
y siempre devuelven "sigue". Los condactos de **control** afectan al ciclo
del intérprete (§7.2) — señalan al motor qué tablas disparar después de la
respuesta actual, o terminan el proceso/turno.

`NOT` no es aplicable a ninguno (todos son acciones, no condiciones).

#### Salida de texto

| Opcode | Mnemónico  | Firma                              | Semántica                                       |
|--------|------------|------------------------------------|-------------------------------------------------|
| `0xC0` | `WRITE`    | `<MsgID> [<Num\|VarID>...]`        | Imprime el mensaje sin salto de línea final.    |
| `0xC1` | `WRITELN`  | `<MsgID> [<Num\|VarID>...]`        | Como `WRITE` y añade `\n` al final.             |
| `0xC2` | `NEWLINE`  | *(sin args)*                       | Imprime un `\n`. Útil para encadenar `WRITE`s.  |
| `0xC3` | `CLEAR`    | *(sin args)*                       | Limpia la pantalla.                             |
| `0xC4` | `LISTOBJ`  | *(sin args)*                       | Imprime la lista de items visibles en la localidad actual, usando el mensaje configurado en `SECTION MESSAGES` (e.g. `loc-objects: "Aquí hay _"`). |
| `0xC5` | `INVENT`   | *(sin args)*                       | Imprime la lista de items que lleva el humano (atajo equivalente a iterar). |

**Interpolación en `WRITE`/`WRITELN`.** El mensaje referenciado puede contener
placeholders `_` (mecanismo ya soportado por el subsistema `Message`,
ver `message.go:CountPlaceholders` y `Stringf`). La firma del condacto
requiere **un valor por cada `_`**, en orden de aparición:

- Cada valor puede ser un literal `Num` (`42`, `-3`) o una `VarID` (label de
  una variable, cuyo valor se resuelve en runtime).
- El compilador **cuenta los `_` del mensaje** y exige exactamente esa
  cantidad de argumentos (§9.1: aridad incorrecta es error de compilación).
- Para mensajes plurales (con `.zero`/`.one`/`.many` en el fuente), el
  **primer** valor dirige la elección de la variante, igual que hace
  `Message.pluralize` hoy.

Ejemplos:

```
// Mensaje no parametrizado:
WRITE seen-cave

// Mensaje con placeholders:
//   you-have-points: "Tienes _ puntos."
WRITE you-have-points score

// Mensaje plural (sección messages):
//   order_count.zero: "No has dado ninguna orden."
//   order_count.one:  "Has dado una orden."
//   order_count.many: "Has dado _ órdenes."
WRITE order_count turns
```

Aridad incorrecta:

```
WRITE you-have-points              // error: faltan 1 arg
WRITE seen-cave 5                   // error: sobra 1 arg
```

#### Audio y visual

| Opcode | Mnemónico | Firma         | Semántica                                                |
|--------|-----------|---------------|----------------------------------------------------------|
| `0xD0` | `PLAY`    | `<BlobID>`    | Arranca la reproducción de un blob de audio. Asíncrono — no bloquea el turno. Si ya había audio sonando, se sustituye. |
| `0xD1` | `STOP`    | *(sin args)*  | Detiene la reproducción de audio actual, si la hay. No-op si nada sonaba. |
| `0xD2` | `PICTURE` | `<BlobID>`    | Muestra el blob como ilustración de la localidad actual. Lo que se ve hasta que el siguiente `PICTURE` lo reemplace. |

Validación en compilación: el `BlobID` debe corresponder a un blob con MIME
compatible (audio para `PLAY`, imagen para `PICTURE`). Se chequea contra el
campo `Mime` del store de Blobs (§9.1 amplía).

#### Pausa

| Opcode | Mnemónico | Firma     | Semántica                                                  |
|--------|-----------|-----------|------------------------------------------------------------|
| `0xE0` | `WAIT`    | `<Num>`   | Bloquea el intérprete durante `num` segundos. Durante la pausa **no se procesan crons** — el reloj de cron también se pausa. Pensado para efecto narrativo entre escrituras. |

#### Control del ciclo

Estos condactos no cortan el proceso; **arman flags** que el ciclo principal
del intérprete (§7.2) consulta al terminar `response` para decidir si dispara
`location` y/o `npc`.

| Opcode | Mnemónico | Firma         | Semántica                                                    |
|--------|-----------|---------------|--------------------------------------------------------------|
| `0xF0` | `LOOK`    | *(sin args)*  | Arma `look_requested`: tras la respuesta, se dispara la tabla `location` para la localidad actual. |
| `0xF1` | `ENGAGE`  | `<CharID>`    | Arma `npc_engaged = char`: tras la respuesta, se dispara la tabla `npc` con la cabecera de ese personaje. |

Usar `LOOK` y `ENGAGE` desde la tabla `response` es la forma estándar de
conectar la entrada del jugador con las tablas de descripción y diálogo.
Por ejemplo, un proceso `examine .*:` típicamente acaba en `LOOK` `DONE`.

#### Terminadores

Repetimos aquí (ya enunciados al final de §6.1/§6.2) los tres condactos que
cortan ejecución:

| Opcode | Mnemónico | Firma         | Semántica                                                |
|--------|-----------|---------------|----------------------------------------------------------|
| `0xFA` | `END`     | *(sin args)*  | Termina este proceso. El motor pasa al siguiente proceso de la misma tabla. Equivale a una condición falsa final. |
| `0xFB` | `DONE`    | *(sin args)*  | Termina este turno: ninguna tabla más se ejecuta hasta el próximo disparo (§7.3). |
| `0xFC` | `OK`      | *(sin args)*  | Equivalente a `WRITE ok` (mensaje canónico) + `DONE`. Atajo del patrón "acción consumida". |

#### Resumen §6.3

15 condactos de salida y control:

`WRITE`, `WRITELN`, `NEWLINE`, `CLEAR`, `LISTOBJ`, `INVENT`,
`PLAY`, `STOP`, `PICTURE`,
`WAIT`,
`LOOK`, `ENGAGE`,
`END`, `DONE`, `OK`.

Opcodes asignados: `0xC0..0xC5`, `0xD0..0xD2`, `0xE0`, `0xF0..0xF1`,
`0xFA..0xFC`. Rango `0xC0..0xFF` da hasta 64 opcodes; quedan ~49 huecos para
el futuro (`BEEP`, `PAUSE_INPUT`, `RAMSAVE`, `RAMLOAD`, etc.).

#### Resumen total del catálogo

**46 condactos en el MVP**, distribuidos en:

- **§6.1 Condiciones** (15): `AT`, `HERE`, `CARRIED`, `WORN`, `PRESENT`,
  `EXIST`, `IN_ROOM`, `IS_OPEN`, `EQ`, `NE`, `GT`, `LT`, `ZERO`, `HUMAN`,
  `CHANCE`. Opcodes `0x01..0x0F`.
- **§6.2 Acciones de mundo** (16): `MOVE`, `PUT_AT`, `GIVE`, `PUT_INTO`,
  `GOTO`, `GET`, `DROP`, `WEAR`, `UNWEAR`, `OPEN`, `CLOSE`, `CREATE`,
  `DESTROY`, `SET`, `INC`, `DEC`. Opcodes `0x40..0x82`.
- **§6.3 Salida y control** (15): listados arriba. Opcodes `0xC0..0xFC`.

Cobertura confirmada:

- 7 condiciones de localización/inventario, 5 de variables, 1 de personaje, 1 de aleatorio, 1 de contenedor.
- 6 acciones de movimiento, 4 de inventario, 2 de contenedor, 2 de existencia, 3 de variables.
- 6 acciones de salida de texto, 3 audiovisuales, 1 de pausa, 2 de control de ciclo, 3 terminadores.

Esto cumple holgadamente los criterios del §11 (≥20 condactos con el reparto descrito).

---

## 7. Semántica de ejecución

> **Estado: propuesta**, pendiente de validación.

### 7.1 Estado global del intérprete

El runtime mantiene, además de la base de datos importada y congelada:

- **Contador de turno** (`uint32`): se incrementa una vez por iteración del
  loop principal **después** de procesar la entrada del jugador. `cron` y
  `init` no lo incrementan.
- **Localidad actual del humano** (`LocID`): se consulta para decidir si
  `location` debe dispararse.
- **Localidad previa**: para detectar cambio.
- **Flags de turno**:
  - `npc_engaged: CharID | nil` — si la entrada del jugador inició una
    interacción con un NPC, qué NPC.
  - `look_requested: bool` — si la entrada incluyó "look" (o el verbo
    equivalente en la lengua de la aventura).
  - `turn_done: bool` — si algún condacto ya hizo `DONE`/`OK`.
- **Cola de `cron` pendientes**: lista de procesos cron disparados durante
  una ejecución bloqueada que aún no se han ejecutado.

### 7.2 Ciclo principal

Pseudocódigo:

```
run table init        // una vez, al arrancar
loop:
    drain pending crons          // (ver 7.4)
    input = read_player()
    if input is blank: continue
    sl = parse(input)            // -> verb+noun, o NullWord (_) en su slot
    turn_done = false
    run table response with sl
    if not turn_done:
        if location_changed or look_requested:
            run table location for current_location
        if npc_engaged:
            run table npc for npc_engaged
    turn_counter += 1
    reset turn flags
```

Notas:

- `init` corre **una sola vez**. Si quieres re-iniciar la aventura, no es
  responsabilidad del lenguaje: el runtime arranca de cero.
- `response` siempre se intenta primero ante input no vacío. Si ningún proceso
  casa (ni siquiera el catch-all `* *:`), el motor escribe el mensaje
  `i-dont-understand` (o equivalente configurado) y consume el turno.
- `location` se dispara **sólo si** la localidad cambió durante `response`
  **o** el jugador pidió describirla. Si una respuesta corrió `DONE`/`OK`, no
  se ejecuta location en ese turno.
- `npc` se dispara **sólo si** `response` marcó el flag `npc_engaged` (un
  condacto como `TALK_TO <char>` lo arma). Mismo bloqueo por `DONE`/`OK`.

### 7.3 Política de aborto, repaso

Resumen consolidado (§4.4 + §6.3):

| Condacto / situación   | Efecto                                                    |
|------------------------|-----------------------------------------------------------|
| condición falsa        | corta el proceso actual, sigue siguiente proceso de la tabla |
| `END`                  | corta el proceso actual, sigue siguiente proceso             |
| `DONE`                 | corta tabla y resto de tablas del turno; flag `turn_done`    |
| `OK`                   | escribe mensaje `ok` y aplica `DONE`                          |
| fin de la tabla        | continúa con la siguiente tabla del ciclo                    |

### 7.4 `cron` y mutex

`cron` se modela como **timers independientes**. Cada proceso `EVERY <intervalo>`
tiene su propio reloj que arranca cuando se carga la aventura:

- Para intervalos en `s` / `m` / `XmYs`: reloj de pared.
- Para intervalos sin sufijo: reloj de turnos (decrementa con cada turno
  completado).

Cuando vence un cron mientras el intérprete está **ocupado** (procesando
`response`/`location`/`npc`), el cron se encola en `pending_crons` y **no se
ejecuta concurrentemente**. Al terminar la tabla en curso y antes de leer la
siguiente entrada del jugador, se vuelca la cola en orden de aparición en el
fuente (`drain pending crons`).

**Coalescencia**. Si el mismo proceso cron se acumula varias veces (porque
estuvo bloqueado mucho tiempo), se ejecuta **una sola vez** por vaciado. No
queremos compensar el desfase con ráfagas — preferimos saltar disparos
perdidos. Esto es una elección consciente, distinta de cron(8) tradicional.

**`DONE` dentro de cron**. Si un proceso cron ejecuta `DONE`, sólo termina ese
cron, no afecta al ciclo principal. El siguiente cron pendiente sigue.

### 7.5 Resolución por orden en una tabla

Dentro de una tabla, los procesos se evalúan en **orden de aparición** en el
fuente. Es responsabilidad del programador situar los procesos específicos
antes que los wildcards (§2.6).

Un proceso es **candidato** si su cabecera casa con la situación actual
(según §4.1). Para tablas como `response` con wildcards, varios procesos
pueden ser candidatos: el motor los prueba en orden y avanza al siguiente
sólo si el actual abortó con condición falsa o `END`. `DONE`/`OK` cortan
la búsqueda.

### 7.6 Determinismo

Salvo `cron` (cuya temporización es real) y `CHANCE <p>` (que introduce
aleatoriedad), la ejecución es **determinista**: dado el mismo estado y la
misma entrada, el resultado es idéntico. Esto facilita los tests de
regresión sobre aventuras compiladas.

---

## 8. Almacenamiento

> **Estado: propuesta**, pendiente de validación.

### 8.1 Decisión: bytecode por proceso (no condactos como registros)

Cada `Process` lleva su lista de condactos **serializada como un blob de bytes
dentro del propio registro**, no como registros separados de tipo `Condact`.

Razones:

- Un condacto es **una pseudo-instrucción**, no una entidad referenciable. No
  hay queries del tipo "dame todos los condactos `WRITE`".
- Atomicidad: un proceso es válido o no; no hay estados intermedios donde
  algunos condactos existen y otros no.
- Compacidad y velocidad: el intérprete itera bytes sin lookups en la BD.
- Aventura típica ≈ 100 procesos × 10 condactos = ~1 KB de bytecode total.
  Negligible.

El tipo `Condact` actual (`internal/adventure/process/condact.go`) seguirá
existiendo **a nivel del compilador** como AST intermedio, pero no se
persiste como `Storeable`.

### 8.2 Tipos persistidos

Tres kinds participan:

- `kind.Table` — un registro por tabla declarada.
- `kind.Process` — un registro por proceso (cabecera + bytecode).
- `kind.Word`, `kind.Location`, etc. — ya existen; los condactos sólo guardan
  IDs hacia ellos.

#### `Table`

Campos persistidos (ampliación del struct actual en `table.go`):

```
Table {
    ID         uint32
    LabelID    uint32           // 'init', 'location', etc.
    Kind       TableKind        // Init / Location / Response / Cron / NPC
    ProcessIDs []uint32         // procesos en orden de aparición en el fuente
}
```

`ProcessIDs` preserva el orden del fuente (clave para §7.5).

#### `Process`

Ampliación del struct actual en `process.go`:

```
Process {
    ID         uint32
    LabelID    uint32           // label de la cabecera o 0 si init plana
    TableID    uint32           // tabla a la que pertenece
    Kind       TableKind        // duplicado de Table.Kind para optimizar dispatch
    Header     ProcessHeader    // cabecera tipada (ver 8.3)
    Bytecode   []byte           // los condactos serializados
}
```

El struct `Process` siempre persiste sus dos slots de cabecera como IDs
(`Slot1`, `Slot2 uint32`); el significado de cada slot depende del `Kind`
de la tabla (§4.1). Cuando el slot vale `*` o `_`, se almacenan los IDs
reservados (1 y 2) — no hay un valor "ausente". `init` guarda `(2, 2)` en
sus dos slots (NullWord en ambos).

#### `ProcessHeader`

La estructura es **flat**: dos campos uniformes que el intérprete decodifica
con el `Kind` de la tabla. Esto evita unión etiquetada y simplifica el
binario.

| `Kind`      | Interpretación de `Slot1`                | Interpretación de `Slot2`                |
|-------------|------------------------------------------|------------------------------------------|
| `Init`      | siempre `_` (id=2)                       | siempre `_` (id=2)                       |
| `Location`  | `LocID` o `*` (id=1)                     | siempre `_` (id=2)                       |
| `Response`  | `WordID` (verbo), o `*`/`_`              | `WordID` (nombre), o `*`/`_`             |
| `Cron`      | unidad codificada como uint32 enum (ver más abajo) | cuenta `uint32`                    |
| `NPC`       | `WordID` (nombre del personaje), o `*`   | `WordID` (adjetivo), o `*`/`_`           |

Para `Cron`, `Slot1` toma uno de cuatro valores enteros bien definidos:

```
1 = HOURS
2 = MINUTES
3 = SECONDS
4 = TURNS
```

Los wildcards no se "expanden" en compilación: se guardan literalmente como
IDs de los labels reservados `*` (id=1) y `_` (id=2), y el dispatcher en
runtime los reconoce.

#### Extensiones a kinds existentes

El catálogo de §6 introduce un campo nuevo en un struct ya existente:

**`Item`** (`internal/adventure/item/item.go`) — añadir:

```
Open  bool   // sólo significativo si Container == true; cerrado por defecto
```

Semántica del campo y los condactos que lo manipulan ya están especificados
(`IS_OPEN` en §6.1, `OPEN`/`CLOSE` en §6.2). Al ser un campo nuevo,
las DBs antiguas no lo traen y deben rehidratarse con `Open = false` por
defecto (cbor ignora campos ausentes — sin migración explícita necesaria).

No se introducen campos nuevos en `Character`, `Location`, `Message`,
`Word`, `Variable`, `Config` ni `Blob`. El resto de §6 se apoya en los
flags ya existentes (`Item.Worn`, `Item.Created`, `Character.Human`,
`Item.At`, `Variable.Value`, etc.).

### 8.3 Formato del bytecode

Cada condacto en el bytecode es:

```
[ opcode_byte | flags_byte | param1_type_byte | param1_payload | ... ]
```

- **opcode** (1 byte): identifica el condacto. Sin valores reservados:
  los 256 espacios están disponibles (la longitud del blob se conoce por
  prefijo, ver más abajo).
- **flags** (1 byte): bit 0 = `NOT` (negación, §4.5). Resto reservado.
- **params**: aridad y tipos fijos por opcode (firma conocida en compilación
  → no hace falta longitud explícita por parámetro). Cada param se serializa
  como:
  - `Str`: `uint16 length + length bytes` (UTF-8).
  - `Num`: `int32`.
  - `Bool`: `uint8` (0/1).
  - `*ID` (`WordID`/`LocID`/…/`CharID`/`BlobID`): `uint32`.
- Endianness: **little-endian**, consistente con `internal/database/io_export.go`.

No se serializa el tipo del param porque la firma del opcode lo determina.

**Delimitación del bytecode**. El blob completo de un `Process` lleva un
prefijo `uint32` con la longitud en bytes del bytecode, igual que el
`Data []byte` de cada `Record` se serializa con `uint64 len` en
`io_export.go`. El intérprete itera hasta agotar esa longitud; no hay
sentinel especial. Esto deja los 256 opcodes disponibles y es más robusto
a corrupciones (un parámetro cero accidental no se confunde con fin de
proceso).

### 8.4 Integridad referencial

Todos los `*ID` se resuelven **en compilación**. El bytecode contiene `uint32`s
listos para usar contra los stores existentes (`Locations.Get().WithID(...)`,
etc.) o las labels reservadas en el caso de wildcards.

En runtime la BD está congelada (§adventure.Import) y los IDs son estables:
ningún `Create` posterior los desplaza (gracias al fix del bug 10 que avanza
`lastDataID`/`lastLabelID` tras Import).

Las **mutaciones de estado** (`Update` de items, characters, variables) pasan
por el sistema de snapshots de `database`. El bytecode no muta — sólo lee y
escribe registros referenciados, no se reescribe a sí mismo.

### 8.5 Versionado

El stream binario que produce `Export` ya incluye un byte `version` por
fichero (`io_export.go:11`). Cuando se añadan nuevos condactos o cambie el
layout, se incrementa `exportVersion` y `Import` decide si puede leer la DB
(actualmente: warning si versión distinta, sigue cargando).

Para el bytecode interno, no hay versión separada: va atado a la versión
global de la BD.

### 8.6 Resumen visual

```
Adventure.Tables.Get().WithLabel("response").First()
  ↓
Table {Kind: Response, ProcessIDs: [10, 11, 12, 13]}
  ↓                                  ↑
  Adventure.Processes.Get().WithID(10)
  ↓
Process {
    Kind: Response,
    Header: {VerbID: take_id, NounID: sword_id},
    Bytecode: <HERE 0x06 sword_id> <NOT 0x01 CARRIED 0x05 sword_id> <GET ...>
}
```

---

## 9. Errores

> **Estado: propuesta**, pendiente de validación.

Los errores se reparten en dos clases. La política rectora: **detectar todo lo
posible en compilación; el runtime sólo falla por imposibilidades genuinas**.

### 9.1 Errores de compilación

Detectados por el compilador y reportados con `CompilerError`
(`internal/compiler/compiler_error/`):

**Estructurales (§2 / §3 / §4):**

- `END TABLE` sin `TABLE <label>` previo, o `TABLE` sin `END TABLE`.
- `TABLE <label>` con `<label>` que no es una de las cinco reservadas.
- Etiqueta de tabla duplicada (`init` declarada dos veces).
- Cabecera de proceso fuera del cuerpo de una tabla.
- Cabecera de `cron` con un intervalo malformado (`EVERY` sin valor, `EVERY 5x`, etc.).
- `EVERY` o cabecera de tipo equivocado en la tabla equivocada (p. ej. `EVERY 5s:` en `response`).

**Firmas y parámetros (§5.4):**

- Aridad incorrecta del condacto.
- Tipo primitivo incompatible.
- Label no resuelta.
- Kind incorrecto (`AT sword` cuando `sword` es un Item).
- Opcode (mnemónico) desconocido.
- `NOT` aplicado a un condacto que no es condición.
- Comodín (`*` / `_`) usado como parámetro de condacto en vez de en cabecera.

**Semántica de proceso:**

- `response`: cabecera con palabras que no son `Verb` + `Noun` (un adjetivo en
  posición de verbo, por ejemplo).
- `location` / `npc`: cabecera con label que no corresponde al kind esperado
  (e.g. `npc` con un label de localidad).
- Proceso `cron` cuya cabecera repite el mismo intervalo de otro proceso de la
  misma tabla → **warning**, no error: es legal pero sospechoso.

**Validación a nivel sección** (en `validateSection` al cerrar la sección):

- Si la aventura define alguna tabla pero falta `init`, **warning**: la
  aventura puede no inicializar variables/posición y comportarse raro.
  Mi recomendación de diseño es **no obligar `init`** — algunas aventuras
  pueden no necesitarla.
- Si `cron` define un `EVERY` con intervalo cero (`EVERY 0s`, `EVERY 0`) →
  error.

### 9.2 Errores de runtime

Detectados durante la interpretación del bytecode. Política:

- Por defecto, **un error de runtime aborta el proceso actual** (como una
  condición falsa) y se registra en el log con `log.Warning`.
- El motor **no entra en pánico**: una aventura mal portada degrada elegante,
  no tira el juego.
- En modo `--strict` (CLI flag, fuera del alcance del MVP) un error de runtime
  abortaría el juego para facilitar debug.

Errores típicos:

- Referencia rota: el bytecode apunta a un `ItemID` que ya no existe
  (sólo puede pasar si la BD se ha corrompido — IDs son estables tras Import).
- División por cero en un `LET <var> <expr>` (cuando se introduzca aritmética).
- Desbordamiento de variable (un `Num` que excede `int32`).
- `cron` con intervalo de turnos = 0 que escapó a la validación de §9.1.
- Bytecode corrupto (longitud inconsistente con el contenido leído).

### 9.3 Reporte y diagnóstico

En compilación: cada error apunta al fichero, número de línea y stack de
contexto (gracias a la infraestructura existente de `CompilerError`).

En runtime: cada error de proceso se registra con:

- Identificación del proceso (`Table.Label + Process.Label` o índice).
- Opcode que falló.
- Estado relevante del juego.

No se imprime nada al jugador salvo que el programador haya escrito
explícitamente `WRITE` o equivalente.

---

## 10. Extensibilidad

> **Estado: propuesta**, pendiente de validación.

### 10.1 Patrón de registro

Cada condacto se define en un único sitio como una **entrada en un registro
estático en Go**. El patrón es paralelo al que ya usa `pkg/validator` con su
slice `validators` (`pkg/validator/validator.go:15`).

Esbozo del struct:

```go
type Condact struct {
    Opcode     byte             // identidad estable en el bytecode
    Mnemonic   string           // como aparece en el fuente
    Aliases    []string         // mnemónicos alternativos opcionales
    Signature  []param.Type     // tipos esperados, en orden
    IsCondition bool            // true → puede ser negada con NOT
    AllowedIn  TableSet         // tablas donde está permitido (vacío = todas)
    Handler    HandlerFunc      // ejecutor en runtime
    Validator  ValidatorFunc    // chequeo extra en compilación (opcional)
}

type HandlerFunc func(state *runtime.State, params []param.Value) Result
type ValidatorFunc func(line.Line, params []param.Value) error
```

El **único sitio** donde se enumeran los condactos es un slice global
`condacts.All` (similar a `validator.validators`). El compilador y el
intérprete leen de ahí: añadir un condacto = añadir una línea al slice +
escribir su `Handler` y, si procede, su `Validator`.

### 10.2 Garantías de identidad

- El `Opcode` (byte) es **estable** entre versiones. Cambiarlo rompe la
  compatibilidad de DBs ya exportadas; sólo se hace al subir
  `exportVersion`.
- Los mnemónicos pueden cambiar libremente (afectan al fuente `.adv`, no
  al binario).
- Añadir un condacto nuevo: usa un opcode libre, no se incrementa la versión.
- Eliminar o renombrar un mnemónico: posible pero rompe fuentes previos. El
  opcode subyacente puede seguir existiendo "huérfano" durante una transición.

### 10.3 Cómo añadir un condacto

Pasos rituales para añadir, por ejemplo, `BEEP`:

1. Definir el opcode libre en `internal/adventure/process/condact_table.go`.
2. Escribir el `Handler` en `internal/runtime/condact/` (módulo a crear con
   la implementación; cada condacto en un fichero o agrupado por tema).
3. Añadir la entrada al slice global con su mnemónico, firma y handler.
4. Añadir tests en `condact_test.go` (compilación: fuente válido / inválido;
   runtime: estado antes / estado después).

No hace falta tocar parser, validador ni intérprete: ambos consultan el
registro.

### 10.4 Plugins externos

**Fuera del MVP.** El registro es estático en Go: no se cargan condactos
desde ficheros, libs dinámicas ni scripts. Razones:

- Seguridad: una aventura no puede inyectar código.
- Determinismo: el set de opcodes es conocido en build-time.
- Simplicidad del runtime: no hay loader.

Si en el futuro se quisiera scripting, sería un mecanismo distinto
(p. ej. un opcode `CALL_SCRIPT <id>` con sandbox) — no se mezcla con
condactos.

### 10.5 Política de catálogo MVP

El catálogo de §6 nace **pequeño y completo**: cubre lo necesario para una
aventura tipo "Aventura Original" y no más. La extensión natural (más
condactos numéricos, manipulación de strings, audio, gráficos avanzados) se
hace a posteriori siguiendo §10.3.

---

## 11. Fuera de alcance

Pendiente de ampliar a medida que se cierren los demás bloques. Lista inicial:

- Debugger interactivo y breakpoints.
- Traza de ejecución persistente.
- Scripts externos / hot-reload de procesos.
- Saltos arbitrarios entre procesos (`GOTO` a etiqueta de proceso).
- Integración con motores gráficos o de audio más allá de `BLOB`.
